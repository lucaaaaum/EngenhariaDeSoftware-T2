package user

import (
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
)

// CreateUserCommand são os dados necessários para criar um usuário.
// Recebidos via JSON na requisição POST /users.
type CreateUserCommand struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserCommand são os dados para atualizar um usuário.
// Recebidos via JSON na requisição PUT /users/{id}.
type UpdateUserCommand struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// UserDto é o que retornamos ao cliente — nunca inclui senha ou hash.
type UserDto struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func NewUserDto(u *user.User) *UserDto {
	return &UserDto{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}
