package kafka

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	producer := NewProducer()
	assert.NotNil(t, producer)
	producer.Close()
}

func TestProduce(t *testing.T) {
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	producer := NewProducer()
	defer producer.Close()

	partion, offset, err := producer.Produce("test_topic", "test_key", "test_message")
	assert.Nil(t, err)
	assert.NotEqual(t, -1, partion)
	assert.NotEqual(t, int64(-1), offset)
}

func TestNewConsumer(t *testing.T) {
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	consumer := NewConsumer("test_group", []string{"test_topic"})
	assert.NotNil(t, consumer)
	consumer.consumer.Close()
}

func TestConsume(t *testing.T) {
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	os.Setenv("TIMEOUT_MS", "5000")
	producer := NewProducer()
	//create a random string for the topic consumer
	topic := "test_topic_" + time.Now().Format("20060102150405.000")
	//create a random consumer group for the topic
	Consumer_group := "test_group_" + time.Now().Format("20060102150405.000")
	// Produce a message
	log.Println("Producing message to topic:", topic)
	_, _, err := producer.Produce(topic, "", "test_message")
	assert.Nil(t, err)

	log.Println("Sleeping for 3 seconds")
	// Allow some time for the message to be produced and consumed
	time.Sleep(3 * time.Second)
	defer producer.Close()

	consumer := NewConsumer(Consumer_group, []string{topic})
	defer consumer.consumer.Close()

	// Consume the message
	message, err := consumer.Consume("")
	assert.Nil(t, err)
	assert.Equal(t, "test_message", message)
}
