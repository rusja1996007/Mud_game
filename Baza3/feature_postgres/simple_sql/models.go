package simple_sql

import "time"

// С Больших букв чтобы были публичными

type TaskModel struct {
	ID           int
	Title        string
	Description  string
	Completed    bool
	Created_at   time.Time
	Completed_at *time.Time
}
