package entities

import "time"

type Todo struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	DueDate       *time.Time `json:"due_date"`       //ポインタ型にして空フォームを送られた時にnilにする
	CompletedDate *time.Time `json:"completed_date"` //ポインタ型にして空フォームを送られた時にnilにする
	CreatedAt     time.Time  `json:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at"` //ポインタ型にして空フォームを送られた時にnilにする
}

type Todos []Todo
