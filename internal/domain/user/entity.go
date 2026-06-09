package user

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generating user id: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	return &User{
		Id:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
