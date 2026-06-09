package handler

import (
	"errors"
	"net/http"
	"strings"
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
		return nil, err
	}
	u, err := h.service.CreateUser(c.Context(), cmd)
	if err != nil {
		return nil, err
	}
	return user.NewUserDto(u), nil
}

func (h *UserHandler) GetUserById(c fuego.ContextNoBody) (*user.UserDto, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	u, err := h.service.GetUserById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fuego.HTTPError{Status: http.StatusNotFound, Detail: "user not found"}
		}
		return nil, err
	}
	return user.NewUserDto(u), nil
}

func (h *UserHandler) UpdateUser(c fuego.ContextWithBody[user.UpdateUserCommand]) (any, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	cmd.Id = id
	if err := h.service.UpdateUser(c.Context(), cmd); err != nil {
		return nil, err
	}
	c.SetStatus(204)
	return nil, nil
}

func (h *UserHandler) DeleteUser(c fuego.ContextNoBody) (any, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		return nil, err
	}
	c.SetStatus(204)
	return nil, nil
}
