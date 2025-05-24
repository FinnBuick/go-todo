package main

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"text"`
	Done        bool   `json:"completed"`
}

func NewTask(text string, nextID int) Task {
	return Task{
		ID:          nextID,
		Description: text,
		Done:        false,
	}
}

func (t *Task) Toggle() {
	temp := t.Done
	t.Done = !temp
}
