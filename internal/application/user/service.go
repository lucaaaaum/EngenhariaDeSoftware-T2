package user

import (
	"context"
	"errors"
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
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get user by id"), err)
	}
	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, command CreateUserCommand) (*user.User, error) {
	user, err := user.NewUser(command.Name)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to create user"), err)
	}

	err = s.repo.AddUser(ctx, user)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to add user to repository"), err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, command UpdateUserCommand) error {
	user, err := s.repo.GetUserById(ctx, command.Id)
	if err != nil {
		return errors.Join(errors.New("Failed to get user by id"), err)
	}

	user.Name = command.Name

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		err = errors.Join(errors.New("Failed to update user"), err)
	}

	return err
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		err = errors.Join(errors.New("Failed to delete user"), err)
	}
	return err
}
