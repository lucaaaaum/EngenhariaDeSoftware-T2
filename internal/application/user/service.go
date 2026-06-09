package user

import (
	"context"
	"fmt"
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	u, err := user.NewUser(cmd.Name, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddUser(ctx, u); err != nil {
		return nil, fmt.Errorf("saving user: %w", err)
	}

	return u, nil
}

func (s *Service) UpdateUser(ctx context.Context, cmd UpdateUserCommand) error {
	u, err := s.repo.GetUserById(ctx, cmd.Id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	u.Name = cmd.Name

	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
