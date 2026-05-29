package db

import (
	"context"
	"errors"
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) user.Repository {
	return &userRepo{db: db}
}

func (r *userRepo) GetUserById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get user by id"), err)
	}
	return &u, nil
}

func (r *userRepo) AddUser(ctx context.Context, user *user.User) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (id, name) VALUES ($1, $2)",
		user.Id,
		user.Name,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to add user"), err)
	}
	return err
}

func (r *userRepo) UpdateUser(ctx context.Context, user *user.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE users
		SET name = $2
		WHERE id = $1
		`,
		user.Id,
		user.Name,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to update user"), err)
	}
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM users
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		err = errors.Join(errors.New("Failed to delete user"), err)
	}
	return err
}
