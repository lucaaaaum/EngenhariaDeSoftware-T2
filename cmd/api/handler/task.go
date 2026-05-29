package handler

import (
	"errors"
	"tarefas/internal/application/task"

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
		return nil, errors.Join(errors.New("Failed to parse request body"), err)
	}
	createdTask, err := h.service.CreateTask(c.Context(), cmd)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to create task"), err)
	}
	return task.NewTaskDto(createdTask), nil
}

func (h *TaskHandler) GetTaskById(c fuego.ContextNoBody) (*task.TaskDto, error) {
	stringId := c.PathParam("id")
	taskId, err := uuid.Parse(stringId)
	if err != nil {
		return nil, errors.Join(errors.New("Invalid task ID"), err)
	}
	taskFound, err := h.service.GetTaskById(c.Context(), taskId)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to get task"), err)
	}
	return task.NewTaskDto(taskFound), nil
}

func (h *TaskHandler) QueryTasks(c fuego.ContextNoBody) ([]task.TaskDto, error) {
	var err error

	assignedToString := c.QueryParam("assignedTo")

	var assignedTo uuid.UUID
	if assignedToString != "" {
		assignedTo, err = uuid.Parse(assignedToString)
		if err != nil {
			return nil, errors.Join(errors.New("Invalid assignedTo user ID"), err)
		}
	}

	createdByString := c.QueryParam("createdBy")
	var createdBy uuid.UUID
	if createdByString != "" {
		createdBy, err = uuid.Parse(createdByString)
		if err != nil {
			return nil, errors.Join(errors.New("Invalid createdBy user ID"), err)
		}
	}

	tasksFound, err := h.service.QueryTasks(c.Context(), createdBy, assignedTo)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to query tasks"), err)
	}

	taskDtos := make([]task.TaskDto, len(tasksFound))
	for i, t := range tasksFound {
		taskDtos[i] = *task.NewTaskDto(t)
	}

	return taskDtos, nil
}

func (h *TaskHandler) UpdateTask(c fuego.ContextWithBody[task.UpdateTaskCommand]) (any, error) {
	cmd, err := c.Body()
	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse request body"), err)
	}
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.Join(errors.New("Invalid task ID"), err)
	}
	cmd.Id = id
	err = h.service.UpdateTask(c.Context(), cmd)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to update task"), err)
	}
	c.SetStatus(204)
	return nil, nil
}

func (h *TaskHandler) DeleteTask(c fuego.ContextNoBody) (any, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return nil, errors.Join(errors.New("Invalid task ID"), err)
	}
	err = h.service.DeleteTask(c.Context(), id)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to delete task"), err)
	}
	c.SetStatus(204)
	return nil, nil
}
