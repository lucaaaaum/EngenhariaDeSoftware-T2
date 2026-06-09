package task

// TaskPriority representa a prioridade de uma tarefa.
// Usamos int para armazenar no banco de forma eficiente.
type TaskPriority int

const (
	LowPriority    TaskPriority = iota // 0 = baixa
	MediumPriority                     // 1 = média
	HighPriority                       // 2 = alta
)
