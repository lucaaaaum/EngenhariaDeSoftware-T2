package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TaskFilter agrupa todos os filtros possíveis para busca de tarefas.
// Ponteiros (*) significam que o campo é opcional — nil = não filtrar por esse campo.
// Exemplo: GET /tasks?status=1&priority=2&dueBefore=2025-12-31
type TaskFilter struct {
	CreatedBy  uuid.UUID
	AssignedTo uuid.UUID
	Status     *TaskStatus   // nil = todos os status
	Priority   *TaskPriority // nil = todas as prioridades
	DueBefore  *time.Time    // nil = sem filtro de prazo
}

// Repository define as operações de banco de dados para tarefas.
type Repository interface {
	GetTaskById(ctx context.Context, id uuid.UUID) (*Task, error)
	QueryTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)
	AddTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
