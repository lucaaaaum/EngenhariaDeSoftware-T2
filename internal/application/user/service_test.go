package user_test

// Testes da camada de serviço de usuários.
//
// 📌 O QUE É UM TESTE?
// Um teste verifica se uma função faz o que deveria fazer.
// Formato padrão em Go: func TestNomeDoQueTesta(t *testing.T)
//
// 📌 O QUE É UM MOCK?
// Um mock é uma versão falsa de uma dependência (ex: banco de dados)
// que usamos nos testes para não precisar de banco real.
//
// Para rodar: go test ./internal/application/user/...

import (
	"context"
	"testing"
	"tarefas/internal/application/user"

	"github.com/google/uuid"
)

// TestCreateUser_Success testa a criação bem-sucedida de um usuário.
func TestCreateUser_Success(t *testing.T) {
	// Arrange (preparar) — cria o serviço com repositório falso
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	cmd := user.CreateUserCommand{
		Name:     "Eduardo",
		Email:    "edu@teste.com",
		Password: "senha123",
	}

	// Act (agir) — executa o que queremos testar
	result, err := service.CreateUser(ctx, cmd)

	// Assert (verificar) — confere se o resultado é o esperado
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected user, got nil")
	}
	if result.Name != cmd.Name {
		t.Errorf("Expected name %q, got %q", cmd.Name, result.Name)
	}
	if result.Email != cmd.Email {
		t.Errorf("Expected email %q, got %q", cmd.Email, result.Email)
	}
	if result.Id == uuid.Nil {
		t.Error("Expected non-nil UUID, got nil UUID")
	}
}

// TestCreateUser_EmptyName testa que criar usuário sem nome retorna erro.
func TestCreateUser_EmptyName(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	cmd := user.CreateUserCommand{
		Name:     "", // nome vazio — deve falhar
		Email:    "edu@teste.com",
		Password: "senha123",
	}

	_, err := service.CreateUser(ctx, cmd)

	if err == nil {
		t.Fatal("Expected error for empty name, got nil")
	}
}

// TestCreateUser_ShortPassword testa que senha curta retorna erro.
func TestCreateUser_ShortPassword(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	cmd := user.CreateUserCommand{
		Name:     "Eduardo",
		Email:    "edu@teste.com",
		Password: "123", // menos de 6 caracteres — deve falhar
	}

	_, err := service.CreateUser(ctx, cmd)

	if err == nil {
		t.Fatal("Expected error for short password, got nil")
	}
}

// TestGetUserById_Success testa buscar usuário pelo ID.
func TestGetUserById_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	// Primeiro cria o usuário
	created, err := service.CreateUser(ctx, user.CreateUserCommand{
		Name: "Maria", Email: "maria@teste.com", Password: "senha456",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Agora busca pelo ID
	found, err := service.GetUserById(ctx, created.Id)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if found.Id != created.Id {
		t.Errorf("Expected ID %v, got %v", created.Id, found.Id)
	}
}

// TestGetUserById_NotFound testa que buscar ID inexistente retorna erro.
func TestGetUserById_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	fakeId := uuid.New()
	_, err := service.GetUserById(ctx, fakeId)

	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

// TestUpdateUser_Success testa atualização do nome do usuário.
func TestUpdateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	created, _ := service.CreateUser(ctx, user.CreateUserCommand{
		Name: "Carlos", Email: "carlos@teste.com", Password: "senha789",
	})

	err := service.UpdateUser(ctx, user.UpdateUserCommand{
		Id:   created.Id,
		Name: "Carlos Atualizado",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verifica que o nome foi atualizado
	updated, _ := service.GetUserById(ctx, created.Id)
	if updated.Name != "Carlos Atualizado" {
		t.Errorf("Expected updated name, got %q", updated.Name)
	}
}

// TestDeleteUser_Success testa remoção de usuário.
func TestDeleteUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := user.NewService(repo)
	ctx := context.Background()

	created, _ := service.CreateUser(ctx, user.CreateUserCommand{
		Name: "Ana", Email: "ana@teste.com", Password: "senha000",
	})

	err := service.DeleteUser(ctx, created.Id)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Usuário não deve mais existir
	_, err = service.GetUserById(ctx, created.Id)
	if err == nil {
		t.Fatal("Expected error after deletion, got nil")
	}
}
