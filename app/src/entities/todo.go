package entities

type Todo struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	DueDate       string `json:"due_date"`
	CompletedDate string `json:"completed_date"`
	CreatedAt     string `json:"created_at"`
	DeletedAt     string `json:"deleted_at"`
}

type Todos []Todo
