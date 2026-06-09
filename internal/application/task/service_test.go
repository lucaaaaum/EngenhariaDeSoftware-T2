package task_test

// Testes da camada de serviço de tarefas.
// Para rodar: go test ./internal/application/task/...

import (
	"context"
	"testing"

	"tarefas/internal/application/task"
	"tarefas/internal/application/webhook"
	domaintask "tarefas/internal/domain/task"

	"github.com/google/uuid"
)

// noopWebhook cria um serviço de webhook que não faz nada (para testes)
func newTestTaskService() *task.Service {
	return task.NewService(newMockTaskRepo(), webhook.NewService())
}

// TestCreateTask_Success testa criação de tarefa com dados válidos.
func TestCreateTask_Success(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()
	creatorId := uuid.New()

	cmd := task.CreateTaskCommand{
		Title:       "Implementar login",
		Description: "Criar endpoint de autenticação",
		CreatedBy:   creatorId,
		Priority:    domaintask.HighPriority,
	}

	result, err := service.CreateTask(ctx, cmd)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected task, got nil")
	}
	if result.Title != cmd.Title {
		t.Errorf("Expected title %q, got %q", cmd.Title, result.Title)
	}
	if result.Status != domaintask.Pending {
		t.Errorf("Expected status Pending (0), got %d", result.Status)
	}
	if result.Priority != domaintask.HighPriority {
		t.Errorf("Expected priority High (2), got %d", result.Priority)
	}
}

// TestCreateTask_EmptyTitle testa que título vazio retorna erro.
func TestCreateTask_EmptyTitle(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	cmd := task.CreateTaskCommand{
		Title:     "", // título vazio — deve falhar
		CreatedBy: uuid.New(),
	}

	_, err := service.CreateTask(ctx, cmd)

	if err == nil {
		t.Fatal("Expected error for empty title, got nil")
	}
}

// TestGetTaskById_Success testa buscar tarefa pelo ID.
func TestGetTaskById_Success(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	created, err := service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Tarefa teste", CreatedBy: uuid.New(),
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	found, err := service.GetTaskById(ctx, created.Id)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if found.Id != created.Id {
		t.Errorf("Expected ID %v, got %v", created.Id, found.Id)
	}
}

// TestGetTaskById_NotFound testa que ID inexistente retorna erro.
func TestGetTaskById_NotFound(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	_, err := service.GetTaskById(ctx, uuid.New())

	if err == nil {
		t.Fatal("Expected error for non-existent task, got nil")
	}
}

// TestUpdateTask_ChangeStatus testa atualização de status de tarefa.
func TestUpdateTask_ChangeStatus(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	created, _ := service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Tarefa para atualizar", CreatedBy: uuid.New(),
	})

	err := service.UpdateTask(ctx, task.UpdateTaskCommand{
		Id:     created.Id,
		Title:  created.Title,
		Status: domaintask.Completed, // muda para concluída
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	updated, _ := service.GetTaskById(ctx, created.Id)
	if updated.Status != domaintask.Completed {
		t.Errorf("Expected status Completed (2), got %d", updated.Status)
	}
}

// TestUpdateTask_AssignUser testa atribuição de usuário a uma tarefa.
func TestUpdateTask_AssignUser(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()
	userId := uuid.New()

	created, _ := service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Tarefa para atribuir", CreatedBy: uuid.New(),
	})

	err := service.UpdateTask(ctx, task.UpdateTaskCommand{
		Id:         created.Id,
		Title:      created.Title,
		AssignedTo: userId,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	updated, _ := service.GetTaskById(ctx, created.Id)
	if updated.AssignedTo != userId {
		t.Errorf("Expected assigned user %v, got %v", userId, updated.AssignedTo)
	}
}

// TestDeleteTask_Success testa remoção de tarefa.
func TestDeleteTask_Success(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	created, _ := service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Tarefa para deletar", CreatedBy: uuid.New(),
	})

	err := service.DeleteTask(ctx, created.Id)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	_, err = service.GetTaskById(ctx, created.Id)
	if err == nil {
		t.Fatal("Expected error after deletion, got nil")
	}
}

// TestQueryTasks_FilterByStatus testa filtro de tarefas por status.
func TestQueryTasks_FilterByStatus(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()
	creatorId := uuid.New()

	// Cria 2 tarefas e conclui uma
	t1, _ := service.CreateTask(ctx, task.CreateTaskCommand{Title: "Tarefa 1", CreatedBy: creatorId})
	t2, _ := service.CreateTask(ctx, task.CreateTaskCommand{Title: "Tarefa 2", CreatedBy: creatorId})

	service.UpdateTask(ctx, task.UpdateTaskCommand{
		Id: t2.Id, Title: t2.Title, Status: domaintask.Completed,
	})
	_ = t1 // t1 fica como Pending

	// Filtra apenas concluídas
	completed := domaintask.Completed
	tasks, err := service.QueryTasks(ctx, domaintask.TaskFilter{
		CreatedBy: creatorId,
		Status:    &completed,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 completed task, got %d", len(tasks))
	}
	if tasks[0].Status != domaintask.Completed {
		t.Errorf("Expected completed status, got %d", tasks[0].Status)
	}
}

// TestQueryTasks_FilterByPriority testa filtro de tarefas por prioridade.
func TestQueryTasks_FilterByPriority(t *testing.T) {
	service := newTestTaskService()
	ctx := context.Background()

	service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Urgente", CreatedBy: uuid.New(), Priority: domaintask.HighPriority,
	})
	service.CreateTask(ctx, task.CreateTaskCommand{
		Title: "Normal", CreatedBy: uuid.New(), Priority: domaintask.LowPriority,
	})

	high := domaintask.HighPriority
	tasks, err := service.QueryTasks(ctx, domaintask.TaskFilter{Priority: &high})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 high priority task, got %d", len(tasks))
	}
}
