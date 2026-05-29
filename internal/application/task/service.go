package task

import (
	"context"
	"errors"
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
)

type Service struct {
	repo task.Repository
}

func NewService(repo task.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetTaskById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	task, err := s.repo.GetTaskById(ctx, id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get task by id"), err)
	}
	return task, nil
}

func (s *Service) QueryTasks(ctx context.Context, createdBy, assignedTo uuid.UUID) ([]*task.Task, error) {
	tasks, err := s.repo.QueryTasks(ctx, createdBy, assignedTo)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to query tasks"), err)
	}
	return tasks, nil
}

func (s *Service) CreateTask(ctx context.Context, command CreateTaskCommand) (*task.Task, error) {
	task, err := task.NewTask(command.Title, command.Description, command.CreatedBy)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to create task"), err)
	}

	err = s.repo.AddTask(ctx, task)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to add task to repository"), err)
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, command UpdateTaskCommand) error {
	task, err := s.repo.GetTaskById(ctx, command.Id)
	if err != nil {
		return errors.Join(errors.New("Failed to get task by id"), err)
	}

	task.Title = command.Title
	task.Description = command.Description
	task.AssignedTo = command.AssignedTo
	task.Status = command.Status

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		return errors.Join(errors.New("Failed to update task"), err)
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteTask(ctx, id)
	if err != nil {
		return errors.Join(errors.New("Failed to delete task"), err)
	}
	return nil
}
