package user

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func NewUser(name string) (*User, error) {
	if len(name) == 0 {
		return nil, errors.New("User name cannot be empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Join(errors.New("Failed to generate user ID"), err)
	}

	return &User{
		Id:   id,
		Name: name,
	}, nil
}
