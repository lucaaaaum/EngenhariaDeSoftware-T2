package user_test

import (
	"context"
	"errors"
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
)

// mockUserRepo é um repositório falso que usamos nos testes.
// Em vez de conectar ao banco de dados de verdade, ele usa um map em memória.
// Isso torna os testes rápidos e independentes de banco de dados.
type mockUserRepo struct {
	users map[uuid.UUID]*user.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uuid.UUID]*user.User)}
}

func (m *mockUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) AddUser(ctx context.Context, u *user.User) error {
	m.users[u.Id] = u
	return nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, u *user.User) error {
	if _, ok := m.users[u.Id]; !ok {
		return errors.New("user not found")
	}
	m.users[u.Id] = u
	return nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}
