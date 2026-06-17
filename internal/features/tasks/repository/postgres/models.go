package tasks_postgres_repository

import "time"

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	Created_At   time.Time
	Completed_At *time.Time
	AuthorUserID int
}
