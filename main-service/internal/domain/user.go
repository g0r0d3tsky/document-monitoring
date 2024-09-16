package domain

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  []byte
	Role      string
}

func GeneratePasswordHash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
func (u *User) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

func (u *User) Matches(plaintextPassword string) (bool, error) {
	p, err := GeneratePasswordHash(plaintextPassword)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(p, u.Password)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
