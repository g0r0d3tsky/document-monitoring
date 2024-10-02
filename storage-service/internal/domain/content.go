package domain

import (
	"github.com/google/uuid"
	"time"
)

type Content struct {
	ID        uuid.UUID
	Filename  string
	FileType  string
	FileSize  int
	FilePath  string
	Checksum  string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}
