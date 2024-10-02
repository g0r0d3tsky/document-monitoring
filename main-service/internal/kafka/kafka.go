package kafka

import (
	"log/slog"
	"main-service/config"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	sarama.AsyncProducer
}

func New(config *config.Config) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.DefaultVersion
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(config.Kafka.BrokerList, cfg)
	if err != nil {
		slog.Error("create producer", err)
	}
	go func() {
		for err = range producer.Errors() {
			slog.Error("write access log entry", err)
		}
	}()

	return &Producer{producer}, nil
}

func (p Producer) Push(topic string, filename string, message []byte) error {
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(filename),
		Value: sarama.ByteEncoder(message),
	}

	p.Input() <- kafkaMessage

	slog.Info("Message sent to Kafka with topic: %s, filename: %s", topic, filename)

	return nil
}
