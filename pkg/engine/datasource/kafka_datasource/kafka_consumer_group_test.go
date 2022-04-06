package kafka_datasource

import (
	"context"
	"github.com/buger/jsonparser"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	log "github.com/jensneuse/abstractlogger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const defaultPartition = 0

// newMockKafkaBroker creates a MockBroker to test ConsumerGroups.
func newMockKafkaBroker(t *testing.T, topic, group string, fr *sarama.FetchResponse) *sarama.MockBroker {
	mockBroker := sarama.NewMockBroker(t, 0)

	mockMetadataResponse := sarama.NewMockMetadataResponse(t).
		SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
		SetLeader(topic, defaultPartition, mockBroker.BrokerID()).
		SetController(mockBroker.BrokerID())

	mockProducerResponse := sarama.NewMockProduceResponse(t).
		SetError(topic, 0, sarama.ErrNoError).
		SetVersion(2)

	mockOffsetResponse := sarama.NewMockOffsetResponse(t).
		SetOffset(topic, defaultPartition, sarama.OffsetOldest, 0).
		SetOffset(topic, defaultPartition, sarama.OffsetNewest, 1).
		SetVersion(1)

	mockCoordinatorResponse := sarama.NewMockFindCoordinatorResponse(t).
		SetCoordinator(sarama.CoordinatorType(0), group, mockBroker)

	mockJoinGroupResponse := sarama.NewMockJoinGroupResponse(t)

	mockSyncGroupResponse := sarama.NewMockSyncGroupResponse(t).
		SetMemberAssignment(&sarama.ConsumerGroupMemberAssignment{
			Version:  0,
			Topics:   map[string][]int32{topic: {0}},
			UserData: nil,
		})

	mockHeartbeatResponse := sarama.NewMockHeartbeatResponse(t)

	mockOffsetFetchResponse := sarama.NewMockOffsetFetchResponse(t).
		SetOffset(group, topic, defaultPartition, 0, "", sarama.KError(0))

	mockOffsetCommitResponse := sarama.NewMockOffsetCommitResponse(t)
	mockBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest":        mockMetadataResponse,
		"ProduceRequest":         mockProducerResponse,
		"OffsetRequest":          mockOffsetResponse,
		"OffsetFetchRequest":     mockOffsetFetchResponse,
		"FetchRequest":           sarama.NewMockSequence(fr),
		"FindCoordinatorRequest": mockCoordinatorResponse,
		"JoinGroupRequest":       mockJoinGroupResponse,
		"SyncGroupRequest":       mockSyncGroupResponse,
		"HeartbeatRequest":       mockHeartbeatResponse,
		"OffsetCommitRequest":    mockOffsetCommitResponse,
	})

	return mockBroker
}

// testConsumerGroupHandler implements sarama.ConsumerGroupHandler interface for testing purposes.
type testConsumerGroupHandler struct {
	processMessage func(msg *sarama.ConsumerMessage)
	ctx            context.Context
	cancel         context.CancelFunc
}

func newDefaultConsumerGroupHandler(processMessage func(msg *sarama.ConsumerMessage)) *testConsumerGroupHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &testConsumerGroupHandler{
		processMessage: processMessage,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (d *testConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	d.cancel() // ready for consuming
	return nil
}

func (d *testConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (d *testConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		d.processMessage(msg)
		sess.MarkMessage(msg, "") // Commit the message and advance the offset.
	}
	return nil
}

func newTestConsumerGroup(groupID string, brokers []string) (sarama.ConsumerGroup, error) {
	kConfig := mocks.NewTestConfig()
	kConfig.Version = sarama.V2_7_0_0
	kConfig.Consumer.Return.Errors = true
	kConfig.ClientID = "graphql-go-tools-test"
	kConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Start with a client
	client, err := sarama.NewClient(brokers, kConfig)
	if err != nil {
		return nil, err
	}

	// Create a new consumer group
	return sarama.NewConsumerGroupFromClient(groupID, client)
}

func TestKafkaMockBroker(t *testing.T) {
	var (
		testMessageKey   = sarama.StringEncoder("test.message.key")
		testMessageValue = sarama.StringEncoder("test.message.value")
		topic            = "test.topic"
		consumerGroup    = "consumer.group"
	)

	fr := &sarama.FetchResponse{Version: 11}
	mockBroker := newMockKafkaBroker(t, topic, consumerGroup, fr)
	defer mockBroker.Close()

	brokerAddr := []string{mockBroker.Addr()}

	cg, err := newTestConsumerGroup(consumerGroup, brokerAddr)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, cg.Close())
	}()

	called := 0

	// Stop after 15 seconds and return an error.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	processMessage := func(msg *sarama.ConsumerMessage) {
		defer cancel()

		t.Logf("Processed message topic: %s, key: %s, value: %s, ", msg.Topic, msg.Key, msg.Value)
		key, _ := testMessageKey.Encode()
		value, _ := testMessageValue.Encode()
		require.Equal(t, key, msg.Key)
		require.Equal(t, value, msg.Value)
		require.Equal(t, topic, msg.Topic)
		called++
	}

	handler := newDefaultConsumerGroupHandler(processMessage)

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Start consuming. Consume is a blocker call and it runs handler.ConsumeClaim at background.
		errCh <- cg.Consume(ctx, []string{topic}, handler)
	}()

	// Ready for consuming
	<-handler.ctx.Done()

	// Add a message to the topic. KafkaConsumerGroupBridge group will fetch that message and trigger ConsumeClaim method.
	fr.AddMessage(topic, defaultPartition, testMessageKey, testMessageValue, 0)

	// When this context is canceled, the processMessage function has been called and run without any problem.
	<-ctx.Done()

	wg.Wait()

	// KafkaConsumerGroupBridge is stopped here.
	require.NoError(t, <-errCh)
	require.Equal(t, 1, called)
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}

// It's just a simple example of graphql federation gateway server, it's NOT a production ready code.
func logger() log.Logger {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}

	return log.NewZapLogger(logger, log.DebugLevel)
}

func TestKafkaConsumerGroupBridge_Subscribe(t *testing.T) {
	var (
		testMessageKey   = sarama.StringEncoder("test.message.key")
		testMessageValue = sarama.StringEncoder(`{"stock":[{"name":"Trilby","price":293,"inStock":2}]}`)
		topic            = "test.topic"
		consumerGroup    = "consumer.group"
	)

	fr := &sarama.FetchResponse{Version: 11}
	mockBroker := newMockKafkaBroker(t, topic, consumerGroup, fr)
	defer mockBroker.Close()

	// Add a message to the topic. The consumer group will fetch that message and trigger ConsumeClaim method.
	fr.AddMessage(topic, defaultPartition, testMessageKey, testMessageValue, 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cg := NewKafkaConsumerGroupBridge(ctx, logger()) // use abstractlogger.NoopLogger if there is no available logger.

	sc := sarama.NewConfig()
	sc.Version = sarama.V2_7_0_0

	options := GraphQLSubscriptionOptions{
		BrokerAddr:   mockBroker.Addr(),
		Topic:        topic,
		GroupID:      consumerGroup,
		ClientID:     "graphql-go-tools-test",
		saramaConfig: sc,
	}
	next := make(chan []byte)
	err := cg.Subscribe(ctx, options, next)
	require.NoError(t, err)

	msg := <-next
	expectedMsg, err := testMessageValue.Encode()
	require.NoError(t, err)

	value, _, _, err := jsonparser.Get(msg, "data")
	require.NoError(t, err)
	require.Equal(t, expectedMsg, value)
}

func TestKafkaConsumerGroup_StartConsuming_And_Stop(t *testing.T) {
	var (
		testMessageKey   = sarama.StringEncoder("test.message.key")
		testMessageValue = sarama.StringEncoder("test.message.value")
		topic            = "test.topic"
		consumerGroup    = "consumer.group"
	)

	fr := &sarama.FetchResponse{Version: 11}
	mockBroker := newMockKafkaBroker(t, topic, consumerGroup, fr)
	defer mockBroker.Close()

	// Add a message to the topic. The consumer group will fetch that message and trigger ConsumeClaim method.
	fr.AddMessage(topic, defaultPartition, testMessageKey, testMessageValue, 0)

	sc := sarama.NewConfig()
	sc.Version = sarama.V2_7_0_0

	options := GraphQLSubscriptionOptions{
		BrokerAddr:   mockBroker.Addr(),
		Topic:        topic,
		GroupID:      consumerGroup,
		ClientID:     "graphql-go-tools-test",
		saramaConfig: sc,
	}

	cg, err := NewKafkaConsumerGroup(logger(), &options)
	require.NoError(t, err)

	messages := make(chan *sarama.ConsumerMessage)
	cg.StartConsuming(messages)

	msg := <-messages
	expectedKey, _ := testMessageKey.Encode()
	require.Equal(t, expectedKey, msg.Key)

	expectedValue, _ := testMessageValue.Encode()
	require.Equal(t, expectedValue, msg.Value)

	require.NoError(t, cg.Close())

	done := make(chan struct{})
	go func() {
		defer func() {
			close(done)
		}()

		cg.WaitUntilConsumerStop()
	}()

	select {
	case <-time.After(15 * time.Second):
		require.Fail(t, "KafkaConsumerGroup could not closed in 15 seconds")
	case <-done:
	}
}
