package task_test

import (
	"tarefas/internal/domain/task"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateTask_Success(t *testing.T) {
	userId, _ := uuid.NewV7()
	title := "Task Title"
	description := "Task Description"
	priority := task.MediumPriority
	dueDate := time.Now().Add(time.Hour * 10)

	task, err := task.NewTask(title, description, userId, priority, &dueDate)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if task.Title != title {
		t.Fatalf("Expected title %q, got %q", title, task.Title)
	}

	if task.Description != description {
		t.Fatalf("Expected description %q, got %q", description, task.Description)
	}

	if task.Priority != priority {
		t.Fatalf("Expected priority %d, got %d", priority, task.Priority)
	}

	if *task.DueDate != dueDate {
		t.Fatalf("Expected due date %q, got %q", dueDate, task.DueDate)
	}
}

func TestCreateTask_Fail_EmptyTitle(t *testing.T) {
	userId, _ := uuid.NewV7()
	dueDate := time.Now()
	task, err := task.NewTask("", "", userId, task.MediumPriority, &dueDate)

	if task != nil {
		t.Fatalf("Expected task to be nil")
	}

	if err == nil {
		t.Fatalf("Expected to have an error")
	}
}
