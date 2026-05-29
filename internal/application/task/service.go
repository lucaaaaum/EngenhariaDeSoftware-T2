package task

import (
	"context"
	"errors"
	"tarefas/internal/domain/task"
)

type Service struct {
	repo task.Repository
}

func NewService(repo task.Repository) *Service {
	return &Service{repo: repo}
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
