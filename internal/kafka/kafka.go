package kafka

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	// Import the confluent-kafka-go library
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer is a Kafka producer using the confluent-kafka-go library and use string serializer
type Producer struct {
	producer *kafka.Producer
}

// NewProducer creates a new Kafka producer
func NewProducer() *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	return &Producer{producer: producer}
}

// Produce sends a message to a Kafka topic
func (p *Producer) Produce(topic string, key string, message string) (int32, int64, error) {
	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		return -1, -1, err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return -1, -1, m.TopicPartition.Error
	}

	return m.TopicPartition.Partition, int64(m.TopicPartition.Offset), nil
}

// Close closes the Kafka producer
func (p *Producer) Close() {
	p.producer.Close()
}

// Consumer is a Kafka consumer using the confluent-kafka-go library and use string deserializer
type Consumer struct {
	consumer *kafka.Consumer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(groupID string, topics []string) *Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BROKER"),
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	return &Consumer{consumer: consumer}
}

// Consume reads all message from a Kafka topic and return the latest message with a particular key
func (c *Consumer) Consume(key string) (string, error) {
	var last_msg_value string
	timeoutMsStr := os.Getenv("TIMEOUT_MS")
	timeoutMs, err := strconv.Atoi(timeoutMsStr)
	if err != nil {
		log.Fatalf("Invalid TIMEOUT_MS value: %s", timeoutMsStr)
	}

	for {
		msg, err := c.consumer.ReadMessage(time.Millisecond * time.Duration(timeoutMs))
		if msg != nil && msg.Key != nil && string(msg.Key) == key {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			last_msg_value = string(msg.Value)
		}

		if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
			log.Printf("Timed out")
			if last_msg_value == "" {
				return "", errors.New("No message found")
			} else {
				return last_msg_value, nil
			}
		} else if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, err.(kafka.Error).Code())
			return "", err
		}
	}
}
