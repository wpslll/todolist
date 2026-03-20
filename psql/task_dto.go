package psql

import "time"

type TaskDto struct {
	Title       string
	Description string
	IsCompleted bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}
