package task

import (
	"context"
	"fmt"
	"tarefas/internal/application/webhook"
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
)

type Service struct {
	repo    task.Repository
	webhook *webhook.Service
}

func NewService(repo task.Repository, webhookService *webhook.Service) *Service {
	return &Service{repo: repo, webhook: webhookService}
}

func (s *Service) GetTaskById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	t, err := s.repo.GetTaskById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	return t, nil
}

func (s *Service) QueryTasks(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	tasks, err := s.repo.QueryTasks(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) CreateTask(ctx context.Context, command CreateTaskCommand) (*task.Task, error) {
	t, err := task.NewTask(command.Title, command.Description, command.CreatedBy, command.Priority, command.DueDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddTask(ctx, t); err != nil {
		return nil, fmt.Errorf("saving task: %w", err)
	}

	go s.webhook.Notify(webhook.EventTaskCreated, t)

	return t, nil
}

func (s *Service) UpdateTask(ctx context.Context, command UpdateTaskCommand) error {
	t, err := s.repo.GetTaskById(ctx, command.Id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	wasAssigned := t.AssignedTo == uuid.Nil && command.AssignedTo != uuid.Nil
	wasCompleted := t.Status != task.Completed && command.Status == task.Completed

	t.Title = command.Title
	t.Description = command.Description
	t.AssignedTo = command.AssignedTo
	t.Status = command.Status
	t.Priority = command.Priority
	t.DueDate = command.DueDate

	if err := s.repo.UpdateTask(ctx, t); err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	if wasAssigned {
		go s.webhook.Notify(webhook.EventTaskAssigned, t)
	} else if wasCompleted {
		go s.webhook.Notify(webhook.EventTaskCompleted, t)
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
