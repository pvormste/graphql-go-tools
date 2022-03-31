package kafka_datasource

import (
	"context"
	"fmt"
	"github.com/jensneuse/abstractlogger"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConsumerGroup struct {
	log           abstractlogger.Logger
	consumerGroup sarama.ConsumerGroup
	ctx           context.Context
}

type KafkaDataSource struct {
	log      abstractlogger.Logger
	messages chan *sarama.ConsumerMessage
	ctx      context.Context
}

type kafkaConsumerGroupHandler struct {
	messages chan *sarama.ConsumerMessage
	ctx      context.Context
}

func (k *kafkaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Start consuming...")
	return nil
}

func (k *kafkaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Consuming done...")
	close(k.messages)
	return nil
}

func (k *kafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx, cancel := context.WithTimeout(k.ctx, time.Second*5)
		select {
		case k.messages <- msg:
			cancel()
			session.MarkMessage(msg, "") // Commit the message and advance the offset.
		case <-ctx.Done():
			cancel()
		case <-k.ctx.Done():
			cancel()
			return nil
		}
	}
	fmt.Println("consume claim done")
	return nil
}

func (c *KafkaConsumerGroup) newConsumerGroup(options GraphQLSubscriptionOptions) (sarama.ConsumerGroup, error) {
	sc := sarama.NewConfig()
	sc.Version = sarama.V2_7_0_0
	sc.ClientID = options.ClientID
	return sarama.NewConsumerGroup([]string{options.BrokerAddr}, options.GroupID, sc)
}

func (c *KafkaConsumerGroup) startConsuming(ctx context.Context, cg sarama.ConsumerGroup, messages chan *sarama.ConsumerMessage, options GraphQLSubscriptionOptions) {
	handler := &kafkaConsumerGroupHandler{
		messages: messages,
		ctx:      c.ctx,
	}

	go func() {
		select {
		case <-ctx.Done():
		case <-c.ctx.Done():
		}
		if err := cg.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err := cg.Consume(ctx, []string{options.Topic}, handler); err != nil {
		fmt.Println("consume error", err)
		return
	}
}

func (c *KafkaConsumerGroup) Subscribe(ctx context.Context, options GraphQLSubscriptionOptions, next chan<- []byte) error {
	cg, err := c.newConsumerGroup(options)
	if err != nil {
		return err
	}

	messages := make(chan *sarama.ConsumerMessage)
	go c.startConsuming(ctx, cg, messages, options)

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("c.ctx closed")
				return
			case <-ctx.Done():
				fmt.Println("context cancelled")
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}
				// TODO: What about msg.Key and msg.Headers?
				next <- msg.Value
			}
		}
	}()

	return nil
}

var _ GraphQLSubscriptionClient = (*KafkaConsumerGroup)(nil)
