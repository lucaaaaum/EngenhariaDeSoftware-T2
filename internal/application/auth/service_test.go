package auth_test

// Testes do serviço de autenticação.
// Para rodar: go test ./internal/application/auth/...

import (
	"context"
	"errors"
	"os"
	"testing"

	"tarefas/internal/application/auth"
	"tarefas/internal/domain/user"

	"github.com/google/uuid"
)

// mockAuthUserRepo é um repositório falso para os testes de auth.
type mockAuthUserRepo struct {
	users map[string]*user.User // indexado por email
}

func newMockAuthUserRepo() *mockAuthUserRepo {
	return &mockAuthUserRepo{users: make(map[string]*user.User)}
}

func (m *mockAuthUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	for _, u := range m.users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockAuthUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockAuthUserRepo) AddUser(ctx context.Context, u *user.User) error {
	m.users[u.Email] = u
	return nil
}

func (m *mockAuthUserRepo) UpdateUser(ctx context.Context, u *user.User) error {
	m.users[u.Email] = u
	return nil
}

func (m *mockAuthUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	for email, u := range m.users {
		if u.Id == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

// TestLogin_Success testa login com credenciais corretas.
func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	repo := newMockAuthUserRepo()

	// Cria um usuário real (com hash de senha)
	u, err := user.NewUser("Eduardo", "edu@teste.com", "senha123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	repo.AddUser(context.Background(), u)

	service := auth.NewService(repo)

	resp, err := service.Login(context.Background(), auth.LoginCommand{
		Email:    "edu@teste.com",
		Password: "senha123",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.Token == "" {
		t.Error("Expected JWT token, got empty string")
	}
}

// TestLogin_WrongPassword testa login com senha incorreta.
func TestLogin_WrongPassword(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	repo := newMockAuthUserRepo()
	u, _ := user.NewUser("Eduardo", "edu@teste.com", "senha123")
	repo.AddUser(context.Background(), u)

	service := auth.NewService(repo)

	_, err := service.Login(context.Background(), auth.LoginCommand{
		Email:    "edu@teste.com",
		Password: "senhaerrada", // senha incorreta
	})

	if err == nil {
		t.Fatal("Expected error for wrong password, got nil")
	}
}

// TestLogin_UserNotFound testa login com email inexistente.
func TestLogin_UserNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	repo := newMockAuthUserRepo()
	service := auth.NewService(repo)

	_, err := service.Login(context.Background(), auth.LoginCommand{
		Email:    "naoexiste@teste.com",
		Password: "senha123",
	})

	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

// TestValidateToken_Valid testa que um token gerado pode ser validado.
func TestValidateToken_Valid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	repo := newMockAuthUserRepo()
	u, _ := user.NewUser("Eduardo", "edu@teste.com", "senha123")
	repo.AddUser(context.Background(), u)

	service := auth.NewService(repo)
	resp, err := service.Login(context.Background(), auth.LoginCommand{
		Email: "edu@teste.com", Password: "senha123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Valida o token gerado
	userId, err := auth.ValidateToken(resp.Token)

	if err != nil {
		t.Fatalf("Expected valid token, got error: %v", err)
	}
	if userId == "" {
		t.Error("Expected user ID from token, got empty string")
	}
	if userId != u.Id.String() {
		t.Errorf("Expected user ID %v, got %v", u.Id.String(), userId)
	}
}

// TestValidateToken_Invalid testa que token inválido retorna erro.
func TestValidateToken_Invalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	_, err := auth.ValidateToken("token-invalido-qualquer")

	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}
