package db

import (
	"context"
	"fmt"
	"tarefas/internal/domain/task"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) task.Repository {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetTaskById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	var t task.Task
	if err := r.db.GetContext(ctx, &t, "SELECT * FROM tasks WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	return &t, nil
}

func (r *taskRepo) QueryTasks(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error) {
	var tasks []*task.Task
	args := []any{}
	i := 1

	query := "SELECT * FROM tasks WHERE 1=1"

	if filter.CreatedBy != uuid.Nil {
		query += fmt.Sprintf(" AND created_by = $%d", i)
		args = append(args, filter.CreatedBy)
		i++
	}
	if filter.AssignedTo != uuid.Nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", i)
		args = append(args, filter.AssignedTo)
		i++
	}
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, *filter.Status)
		i++
	}
	if filter.Priority != nil {
		query += fmt.Sprintf(" AND priority = $%d", i)
		args = append(args, *filter.Priority)
		i++
	}
	if filter.DueBefore != nil {
		query += fmt.Sprintf(" AND due_date <= $%d", i)
		args = append(args, *filter.DueBefore)
		i++
	}

	query += " ORDER BY created_at DESC"

	if err := r.db.SelectContext(ctx, &tasks, query, args...); err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	return tasks, nil
}

func (r *taskRepo) AddTask(ctx context.Context, t *task.Task) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (id, title, description, status, priority, due_date, created_by, assigned_to)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		t.Id, t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.CreatedBy, t.AssignedTo,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

func (r *taskRepo) UpdateTask(ctx context.Context, t *task.Task) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE tasks
		 SET title = $2, description = $3, status = $4, priority = $5,
		     due_date = $6, created_by = $7, assigned_to = $8
		 WHERE id = $1`,
		t.Id, t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.CreatedBy, t.AssignedTo,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (r *taskRepo) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
