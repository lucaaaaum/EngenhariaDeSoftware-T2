package db

import (
	"context"
	"errors"
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
	err := r.db.GetContext(ctx, &t, "SELECT * FROM tasks WHERE id = $1", id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get task by id"), err)
	}
	return &t, nil
}

func (r *taskRepo) QueryTasks(ctx context.Context, createdBy, assignedTo uuid.UUID) ([]*task.Task, error) {
	var tasks []*task.Task
	var err error
	args := []any{}
	appliedFilters := 1

	query := "SELECT * FROM tasks WHERE 1=1"

	if createdBy != uuid.Nil {
		query += fmt.Sprintf(" AND created_by = $%d", appliedFilters)
		args = append(args, createdBy)
		appliedFilters++
	}

	if assignedTo != uuid.Nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", appliedFilters)
		args = append(args, assignedTo)
		appliedFilters++
	}

	query += " ORDER BY created_at DESC"

	err = r.db.SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to query tasks"), err)
	}
	return tasks, nil
}

func (r *taskRepo) AddTask(ctx context.Context, task *task.Task) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO tasks (id, title, description, status, created_by, assigned_to) VALUES ($1, $2, $3, $4, $5, $6)",
		task.Id,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedBy,
		task.AssignedTo,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to add task"), err)
	}
	return err
}

func (r *taskRepo) UpdateTask(ctx context.Context, task *task.Task) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE tasks
		SET title = $2, description = $3, status = $4, created_by = $5, assigned_to = $6
		WHERE id = $1
		`,
		task.Id,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedBy,
		task.AssignedTo,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to update task"), err)
	}
	return err
}

func (r *taskRepo) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM tasks
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to delete task"), err)
	}
	return err
}
