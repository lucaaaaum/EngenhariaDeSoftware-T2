package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository define as operações de banco de dados para usuários.
// É uma interface — separa a lógica de negócio do banco de dados (Clean Architecture).
type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	AddUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
