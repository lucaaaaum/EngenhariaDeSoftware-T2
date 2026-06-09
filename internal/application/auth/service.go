package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"tarefas/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*LoginResponse, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		// mesmo erro pra email errado e senha errada, pra não vazar info
		return nil, errors.New("invalid email or password")
	}

	if !u.CheckPassword(cmd.Password) {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"sub": u.Id.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret()))
	if err != nil {
		return nil, fmt.Errorf("signing token: %w", err)
	}

	return &LoginResponse{Token: tokenStr}, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret()), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user id not found in token")
	}

	return userId, nil
}

func jwtSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "default-secret-change-in-production"
}
