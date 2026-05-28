package task

type TaskRepository interface {
	GetTaskById(id string) (*Task, error)
	GetTasksAssignedToUser(userId string) ([]*Task, error)
	GetTasksCreatedByUser(userId string) ([]*Task, error)
	AddTask(task *Task) error
	UpdateTask(task *Task) error
	DeleteTask(id string) error
}
