package todo

type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

func (t *Task) Toggle() {
	temp := t.Completed
	t.Completed = !temp
}
