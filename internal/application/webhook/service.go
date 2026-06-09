package webhook

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
	"tarefas/internal/domain/task"
)

type EventType string

const (
	EventTaskCreated   EventType = "task.created"
	EventTaskAssigned  EventType = "task.assigned"
	EventTaskCompleted EventType = "task.completed"
)

type WebhookPayload struct {
	Event     EventType  `json:"event"`
	Timestamp time.Time  `json:"timestamp"`
	Task      *task.Task `json:"task"`
}

type Service struct {
	url    string
	client *http.Client
}

func NewService() *Service {
	return &Service{
		url:    os.Getenv("WEBHOOK_URL"),
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify dispara um POST para o webhook configurado. Deve ser chamado em goroutine.
func (s *Service) Notify(event EventType, t *task.Task) {
	if s.url == "" {
		return
	}

	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now(),
		Task:      t,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("webhook: failed to marshal payload", "error", err)
		return
	}

	resp, err := s.client.Post(s.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("webhook: request failed", "event", event, "error", err)
		return
	}
	defer resp.Body.Close()

	slog.Info("webhook sent", "event", event, "status", resp.StatusCode)
}
