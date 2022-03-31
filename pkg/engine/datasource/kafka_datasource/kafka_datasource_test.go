package kafka_datasource

import (
	"github.com/Shopify/sarama"
	"testing"
)

func TestKafkaDataSource(t *testing.T) {
	var (
		testMessageKey   = sarama.StringEncoder("test.message.key")
		testMessageValue = sarama.StringEncoder("test.message.value")
		topic            = "test.topic"
		consumerGroup    = "consumer.group"
	)

	fr := &sarama.FetchResponse{Version: 11}
	mockBroker := newMockKafkaBroker(t, topic, consumerGroup, fr)
	defer mockBroker.Close()

	// Add a message to the topic. KafkaConsumerGroup group will fetch that message and trigger ConsumeClaim method.
	fr.AddMessage(topic, defaultPartition, testMessageKey, testMessageValue, 0)

}
