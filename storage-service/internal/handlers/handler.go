package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"storage-service/internal/domain"
	"storage-service/internal/kafka"

	"github.com/IBM/sarama"
)

type Service interface {
	CreateURL(ctx context.Context, message *domain.Url) error
}

type ServiceMessage struct {
	service Service
}

func NewServiceMessage(service Service) ServiceMessage {
	return ServiceMessage{
		service: service,
	}
}

func (s *ServiceMessage) Handler() *kafka.Consumer {
	handler := func(msg *sarama.ConsumerMessage) error {
		var url *domain.Url
		err := json.Unmarshal(msg.Value, &url)
		if err != nil {
			slog.Error("unmarshalling json: %v", err)
			return fmt.Errorf("unmarshalling json: %w", err)
		}

		err = s.CreateMessage(context.Background(), url)
		if err != nil {
			slog.Error("creating message: %v", err)
			return fmt.Errorf("creating message: %w", err)
		}

		return nil
	}

	return kafka.NewConsumer(handler)
}

func (s *ServiceMessage) CreateMessage(ctx context.Context, url *domain.Url) error {
	err := s.service.CreateURL(ctx, url)
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}
	return nil
}
