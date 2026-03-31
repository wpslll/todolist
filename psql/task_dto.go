package psql

import "time"

type TaskDto struct {
	Id int
	Title       string
	Description string
	IsCompleted bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}
