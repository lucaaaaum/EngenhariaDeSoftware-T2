package task_test

import (
	"context"
	"errors"
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
)

// mockTaskRepo é um repositório falso de tarefas para testes.
type mockTaskRepo struct {
	tasks map[uuid.UUID]*task.Task
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{tasks: make(map[uuid.UUID]*task.Task)}
}

func (m *mockTaskRepo) GetTaskById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return t, nil
}

func (m *mockTaskRepo) QueryTasks(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	var result []*task.Task
	for _, t := range m.tasks {
		// Aplica filtros
		if filter.CreatedBy != uuid.Nil && t.CreatedBy != filter.CreatedBy {
			continue
		}
		if filter.AssignedTo != uuid.Nil && t.AssignedTo != filter.AssignedTo {
			continue
		}
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *mockTaskRepo) AddTask(ctx context.Context, t *task.Task) error {
	m.tasks[t.Id] = t
	return nil
}

func (m *mockTaskRepo) UpdateTask(ctx context.Context, t *task.Task) error {
	if _, ok := m.tasks[t.Id]; !ok {
		return errors.New("task not found")
	}
	m.tasks[t.Id] = t
	return nil
}

func (m *mockTaskRepo) DeleteTask(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(m.tasks, id)
	return nil
}
