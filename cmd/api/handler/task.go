package handler

import (
	"errors"
	"strconv"
	"time"

	"tarefas/internal/application/task"
	domaintask "tarefas/internal/domain/task"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *task.Service
}

func NewTaskHandler(service *task.Service) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(c fuego.ContextWithBody[task.CreateTaskCommand]) (*task.TaskDto, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, err
	}
	t, err := h.service.CreateTask(c.Context(), cmd)
	if err != nil {
		return nil, err
	}
	return task.NewTaskDto(t), nil
}

func (h *TaskHandler) GetTaskById(c fuego.ContextNoBody) (*task.TaskDto, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid task id")
	}
	t, err := h.service.GetTaskById(c.Context(), id)
	if err != nil {
		return nil, err
	}
	return task.NewTaskDto(t), nil
}

func (h *TaskHandler) QueryTasks(c fuego.ContextNoBody) ([]task.TaskDto, error) {
	filter := domaintask.TaskFilter{}

	if v := c.QueryParam("assignedTo"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return nil, errors.New("invalid assignedTo")
		}
		filter.AssignedTo = id
	}

	if v := c.QueryParam("createdBy"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return nil, errors.New("invalid createdBy")
		}
		filter.CreatedBy = id
	}

	if v := c.QueryParam("status"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("status must be 0, 1 or 2")
		}
		s := domaintask.TaskStatus(n)
		filter.Status = &s
	}

	if v := c.QueryParam("priority"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("priority must be 0, 1 or 2")
		}
		p := domaintask.TaskPriority(n)
		filter.Priority = &p
	}

	if v := c.QueryParam("dueBefore"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			t, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, errors.New("dueBefore must be YYYY-MM-DD")
			}
		}
		filter.DueBefore = &t
	}

	found, err := h.service.QueryTasks(c.Context(), filter)
	if err != nil {
		return nil, err
	}

	result := make([]task.TaskDto, len(found))
	for i, t := range found {
		result[i] = *task.NewTaskDto(t)
	}
	return result, nil
}

func (h *TaskHandler) UpdateTask(c fuego.ContextWithBody[task.UpdateTaskCommand]) (any, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid task id")
	}
	cmd.Id = id
	if err := h.service.UpdateTask(c.Context(), cmd); err != nil {
		return nil, err
	}
	c.SetStatus(204)
	return nil, nil
}

func (h *TaskHandler) DeleteTask(c fuego.ContextNoBody) (any, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.New("invalid task id")
	}
	if err := h.service.DeleteTask(c.Context(), id); err != nil {
		return nil, err
	}
	c.SetStatus(204)
	return nil, nil
}
