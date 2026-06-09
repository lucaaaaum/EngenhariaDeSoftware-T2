package handler

import (
	"net/http"
	"tarefas/internal/application/auth"

	"github.com/go-fuego/fuego"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c fuego.ContextWithBody[auth.LoginCommand]) (*auth.LoginResponse, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, err
	}
	resp, err := h.service.Login(c.Context(), cmd)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusUnauthorized, Detail: "invalid email or password"}
	}
	return resp, nil
}

// Logout é stateless — JWT não tem sessão no servidor, o cliente descarta o token
func (h *AuthHandler) Logout(c fuego.ContextNoBody) (map[string]string, error) {
	return map[string]string{"message": "logged out"}, nil
}
