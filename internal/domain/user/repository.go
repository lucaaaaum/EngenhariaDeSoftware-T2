package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	AddUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
