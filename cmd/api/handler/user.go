package handler

import (
	"errors"
	"tarefas/internal/application/user"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c fuego.ContextWithBody[user.CreateUserCommand]) (*user.UserDto, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse request body"), err)
	}
	createdUser, err := h.service.CreateUser(c.Context(), cmd)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to create user"), err)
	}

	return user.NewUserDto(createdUser), nil
}

func (h *UserHandler) GetUserById(c fuego.ContextNoBody) (*user.UserDto, error) {
	stringId := c.PathParam("id")
	id, err := uuid.Parse(stringId)
	if err != nil {
		return nil, errors.Join(errors.New("Invalid user ID"), err)
	}
	userFound, err := h.service.GetUserById(c.Context(), id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get user"), err)
	}
	return user.NewUserDto(userFound), nil
}

func (h *UserHandler) UpdateUser(c fuego.ContextWithBody[user.UpdateUserCommand]) (any, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse request body"), err)
	}
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.Join(errors.New("Invalid user ID"), err)
	}
	cmd.Id = id
	err = h.service.UpdateUser(c.Context(), cmd)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to update user"), err)
	}
	c.SetStatus(204)
	return nil, nil
}

func (h *UserHandler) DeleteUser(c fuego.ContextNoBody) (any, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.Join(errors.New("Invalid user ID"), err)
	}
	err = h.service.DeleteUser(c.Context(), id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to delete user"), err)
	}
	c.SetStatus(204)
	return nil, nil
}
