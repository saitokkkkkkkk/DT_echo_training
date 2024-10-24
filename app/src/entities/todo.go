package entities

import "time"

type Todo struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	DueDate       *time.Time `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date"`
	CreatedAt     time.Time  `json:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type Todos []Todo
