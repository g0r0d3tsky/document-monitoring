package usecase

import (
	"context"
	"fmt"
	"io/ioutil"
	"main-service/internal/domain"
	"main-service/internal/usecase/model"
	"os"
	"path/filepath"
)

type MetaRepo interface {
	GetContentByName(ctx context.Context, name string) (*domain.Content, error)
	DeleteContent(ctx context.Context, filename string) error
}

type TextContentRepo interface {
	GetContentByFilename(ctx context.Context, filename string) (*model.Content, error)
}

type Broker interface {
	Push(topic string, filename string, message []byte) error
}

type ContentService struct {
	metaRepo   MetaRepo
	textRepo   TextContentRepo
	kafka      Broker
	kafkaTopic string
}

func NewContentService(metaRepo MetaRepo, textRepo TextContentRepo, producer Broker, kafkaTopic string) *ContentService {
	return &ContentService{metaRepo: metaRepo, textRepo: textRepo, kafka: producer, kafkaTopic: kafkaTopic}
}

func (s *ContentService) SendFileToKafka(filename string, fileContent []byte) error {
	err := s.kafka.Push(s.kafkaTopic, filename, fileContent)
	if err != nil {
		return fmt.Errorf("failed to send file to Kafka: %w", err)
	}

	return nil
}

func (s *ContentService) GetFile(ctx context.Context, filename string) ([]byte, error) {
	if isTextFile(filename) {
		content, err := s.textRepo.GetContentByFilename(ctx, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file from MongoDB: %w", err)
		}
		return []byte(content.Text), nil
	} else {
		meta, err := s.metaRepo.GetContentByName(ctx, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file metadata from PostgreSQL: %w", err)
		}

		fileContent, err := fetchFileFromServer(meta.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file from server: %w", err)
		}

		return fileContent, nil
	}
}

func isTextFile(filename string) bool {
	return filepath.Ext(filename) == ".txt"
}

func fetchFileFromServer(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return fileContent, nil
}

func (s *ContentService) DeleteContentFromServer(ctx context.Context, filename string) error {
	filePath := fmt.Sprintf("/path/to/your/files/%s", filename)

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file from server: %w", err)
	}

	err = s.metaRepo.DeleteContent(ctx, filename)
	if err != nil {
		return fmt.Errorf("error deleting content metadata: %w", err)
	}

	return nil
}
