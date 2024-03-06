package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ProducerKafka interface {
	ProduceMessage(topic string, message []byte) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducerInit(brokers []string) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{writer: writer}, nil
}

func (k *KafkaProducer) ProduceMessage(topic string, message []byte) error {
	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: message,
	})
}

func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}
