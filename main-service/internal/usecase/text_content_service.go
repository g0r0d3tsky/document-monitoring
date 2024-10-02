package usecase

import (
	"context"
	"fmt"
	"io"
	"main-service/internal/usecase/model"
	"mime/multipart"
)

type TextRepo interface {
	SaveTextContent(ctx context.Context, fileContent []byte, filename string) error
	GetContentByFilename(ctx context.Context, filename string) (*model.Content, error)
	DeleteContentByFilename(ctx context.Context, filename string) error
}

type TextService struct {
	repo TextRepo
}

func NewTextService(repo TextRepo) *TextService {
	return &TextService{repo: repo}
}

func (uc *TextService) CreateTextContent(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file content: %w", err)
	}

	err = uc.repo.SaveTextContent(ctx, fileContent, header.Filename)
	if err != nil {
		return fmt.Errorf("error saving content: %w", err)
	}

	return nil
}

func (uc *TextService) GetTextContent(ctx context.Context, filename string) (*model.Content, error) {
	content, err := uc.repo.GetContentByFilename(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("error retrieving content: %w", err)
	}

	return content, nil
}

func (s *TextService) DeleteContent(ctx context.Context, filename string) error {
	err := s.repo.DeleteContentByFilename(ctx, filename)
	if err != nil {
		return fmt.Errorf("error deleting content: %w", err)
	}

	return nil
}
