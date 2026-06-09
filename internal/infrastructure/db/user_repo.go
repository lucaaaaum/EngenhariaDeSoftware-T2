package db

import (
	"context"
	"fmt"
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
	if err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *userRepo) AddUser(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)",
		u.Id, u.Name, u.Email, u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *userRepo) UpdateUser(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET name = $2 WHERE id = $1`, u.Id, u.Name)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
