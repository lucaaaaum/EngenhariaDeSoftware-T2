package user

import (
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Name string `json:"name"`
}

type UpdateUserCommand struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUserDto(user *user.User) *UserDto {
	return &UserDto{
		Id:   user.Id,
		Name: user.Name,
	}
}
