package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"main-service/internal/domain"
)

type StorageUser struct {
	db *pgxpool.Pool
}

func NewStorageURL(dbPool *pgxpool.Pool) StorageUser {
	StorageURL := StorageUser{
		db: dbPool,
	}
	return StorageURL
}

func (s *StorageUser) GetUser(ctx context.Context, userName string, password string) (*domain.User, error) {
	user := &domain.User{}
	if err := s.db.QueryRow(
		ctx,
		`SELECT id, userName, firstName, lastName, email, password, role  FROM "users" u WHERE u.userName = $1`, userName,
	).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user, &user.Password, &user.Role); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("get user: %w", err)
	}
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("wrong password: %w", err)
	}

	return user, nil
}

func (s *StorageUser) CreateUser(ctx context.Context, user *domain.User, password string) error {
	hashedPassword, err := domain.GeneratePasswordHash(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error generating uuid: %w", err)
	}

	user.ID = id

	query := `INSERT INTO "users" (id, userName, firstName, lastName, email, password, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err = s.db.Exec(
		ctx,
		query,
		user.ID, user.Username, user.FirstName, user.LastName, user.Email, string(hashedPassword), user.Role,
	); err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}
