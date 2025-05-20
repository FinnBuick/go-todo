package todo

type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

func NewTask(text string, nextID int) Task {
	return Task{
		ID:        nextID,
		Text:      text,
		Completed: false,
	}
}

func (t *Task) Toggle() {
	temp := t.Completed
	t.Completed = !temp
}
