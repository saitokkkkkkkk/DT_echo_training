package entities

import "time"

type Todo struct {
	ID                     int64      `json:"id"`
	Title                  string     `json:"title"`
	Content                string     `json:"content"`
	DueDate                *string    `json:"due_date"`
	CompletedDate          *string    `json:"completed_date"`
	CreatedAt              time.Time  `json:"created_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
	FormattedDueDate       string     `json:"formatted_due_date"`
	FormattedCompletedDate string     `json:"formatted_completed_date"`
}

type Todos []Todo
