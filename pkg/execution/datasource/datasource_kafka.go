package datasource

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/jensneuse/abstractlogger"
	"github.com/jensneuse/graphql-go-tools/internal/pkg/unsafebytes"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer/literal"
)

type KafkaDataSourceConfig struct {
	BrokerAddr string
	ClientID   string
	Topic      string
	GroupID    string
}

type KafkaDataSourcePlannerFactoryFactory struct {
}

type KafkaDataSourcePlannerFactory struct {
	base   BasePlanner
	config KafkaDataSourceConfig
}

type KafkaDataSourcePlanner struct {
	BasePlanner
	dataSourceConfig KafkaDataSourceConfig
}

func (k KafkaDataSourcePlanner) Plan(args []Argument) (DataSource, []Argument) {
	k.Args = append(k.Args, &StaticVariableArgument{
		Name:  literal.BROKERADDR,
		Value: []byte(k.dataSourceConfig.BrokerAddr),
	})
	k.Args = append(k.Args, &StaticVariableArgument{
		Name:  literal.TOPIC,
		Value: []byte(k.dataSourceConfig.Topic),
	})
	k.Args = append(k.Args, &StaticVariableArgument{
		Name:  literal.CLIENTID,
		Value: []byte(k.dataSourceConfig.ClientID),
	})
	k.Args = append(k.Args, &StaticVariableArgument{
		Name:  literal.GROUPID,
		Value: []byte(k.dataSourceConfig.GroupID),
	})
	return &KafkaDataSource{
		Log: k.Log,
	}, append(k.Args, args...)
}

func (k KafkaDataSourcePlannerFactory) DataSourcePlanner() Planner {
	return SimpleDataSourcePlanner(&KafkaDataSourcePlanner{
		BasePlanner:      k.base,
		dataSourceConfig: k.config,
	})
}

func (k KafkaDataSourcePlannerFactoryFactory) Initialize(base BasePlanner, configReader io.Reader) (PlannerFactory, error) {
	factory := &KafkaDataSourcePlannerFactory{
		base: base,
	}
	return factory, json.NewDecoder(configReader).Decode(&factory.config)
}

type KafkaDataSource struct {
	Log      log.Logger
	once     sync.Once
	messages chan *sarama.ConsumerMessage
	ctx      context.Context
	cancel   context.CancelFunc
}

type kafkaConsumerGroupHandler struct {
	messages chan *sarama.ConsumerMessage
}

func (k *kafkaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (k *kafkaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (k *kafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		k.messages <- msg
		session.MarkMessage(msg, "") // Commit the message and advance the offset.
	}
	return nil
}

func (k *KafkaDataSource) newConsumerGroup(args ResolverArgs) (sarama.ConsumerGroup, error) {
	brokerArg := args.ByKey(literal.BROKERADDR)
	clientIDArg := args.ByKey(literal.CLIENTID)
	topicArg := args.ByKey(literal.TOPIC)
	groupIDArg := args.ByKey(literal.GROUPID)
	k.Log.Debug("KafkaDataSource.Resolve.init",
		log.String("broker", string(brokerArg)),
		log.String("clientID", string(clientIDArg)),
		log.String("groupID", string(groupIDArg)),
		log.String("topic", string(topicArg)),
	)

	sc := sarama.NewConfig()
	sc.Version = sarama.V2_7_0_0
	sc.ClientID = string(clientIDArg)
	return sarama.NewConsumerGroup([]string{string(brokerArg)}, string(groupIDArg), sc)
}

func (k *KafkaDataSource) startConsuming(ctx context.Context, cg sarama.ConsumerGroup, args ResolverArgs) {
	topicArg := args.ByKey(literal.TOPIC)
	handler := &kafkaConsumerGroupHandler{messages: k.messages}

	go func() {
		<-ctx.Done()
		if err := cg.Close(); err != nil {
			k.Log.Error("KafkaDataSource.Resolve", log.Error(err))
		}
	}()

	if err := cg.Consume(ctx, []string{unsafebytes.BytesToString(topicArg)}, handler); err != nil {
		k.Log.Error("KafkaDataSource.Resolve",
			log.Error(err),
		)
		k.cancel()
		return
	}
}

func (k *KafkaDataSource) Resolve(ctx context.Context, args ResolverArgs, out io.Writer) (n int, err error) {
	k.once.Do(func() {
		k.ctx, k.cancel = context.WithCancel(context.Background())
		k.messages = make(chan *sarama.ConsumerMessage)

		cg, err := k.newConsumerGroup(args)
		if err != nil {
			k.Log.Error("KafkaDataSource.Resolve",
				log.Error(err),
			)
			k.cancel()
			return
		}
		go k.startConsuming(ctx, cg, args)
	})

	select {
	case <-k.ctx.Done():
		return
	case <-ctx.Done():
		return
	case msg, ok := <-k.messages:
		if !ok {
			return
		}
		// TODO: What about msg.Key and msg.Headers?
		return out.Write(msg.Value)
	}
}
