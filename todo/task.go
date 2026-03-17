package todo

import "time"

type Task struct {
	Title        string
	Description  string
	IsCompleted  bool

	CreatedAt    time.Time
	CompletedAt *time.Time
}

func NewTask(title string, description string) Task {
	return Task {
		Title: title,
		Description: description,
		IsCompleted: false,

		CreatedAt: time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	completedTime := time.Now()

	t.IsCompleted = true
	t.CompletedAt = &completedTime
}

func (t *Task) Uncomplete() {
	t.IsCompleted = false
	t.CompletedAt = nil
}