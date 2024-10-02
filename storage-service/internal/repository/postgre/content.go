package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"storage-service/internal/domain"
)

type StorageContent struct {
	db *pgxpool.Pool
}

func NewStorageContent(dbPool *pgxpool.Pool) StorageContent {
	StorageContent := StorageContent{
		db: dbPool,
	}
	return StorageContent
}

func (s *StorageContent) CreateContent(ctx context.Context, content *domain.Content) error {
	id := uuid.New()
	content.ID = id

	query := `INSERT INTO "files" (id, filename, file_size, file_path, checksum, created_at, updated_at, user_id)
	          VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $6)`

	if _, err := s.db.Exec(
		ctx,
		query,
		content.ID,
		content.Filename,
		content.FileSize,
		content.FilePath,
		content.Checksum,
		content.UserID,
	); err != nil {
		return fmt.Errorf("error inserting file: %w", err)
	}

	return nil
}

func (s *StorageContent) GetContentByName(ctx context.Context, name string) (*domain.Content, error) {

	query := `SELECT id, filename, file_size, file_path, checksum, created_at, updated_at, user_id
	          FROM "files" WHERE filename = $1`

	var content domain.Content

	if err := s.db.QueryRow(ctx, query, name).Scan(
		&content.ID,
		&content.Filename,
		&content.FileSize,
		&content.FilePath,
		&content.Checksum,
		&content.CreatedAt,
		&content.UpdatedAt,
		&content.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("error retrieving file: %w", err)
	}

	return &content, nil
}

func (s *StorageContent) UpdateContent(ctx context.Context, content *domain.Content) error {

	query := `UPDATE "files" 
	          SET filename = $1, file_size = $2, file_path = $3, checksum = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $5`

	if _, err := s.db.Exec(
		ctx,
		query,
		content.Filename,
		content.FileSize,
		content.FilePath,
		content.Checksum,
		content.ID,
	); err != nil {
		return fmt.Errorf("error updating file: %w", err)
	}

	return nil
}

func (s *StorageContent) DeleteContent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "files" WHERE id = $1`

	if _, err := s.db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
