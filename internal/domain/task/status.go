package task

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
)
