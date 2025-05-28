package models

type Task struct {
	ID          int
	Description string
	Done        bool
}

func NewTask(text string, nextID int) Task {
	return Task{
		ID:          nextID,
		Description: text,
		Done:        false,
	}
}
