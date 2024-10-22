package validators

import "github.com/go-playground/validator/v10"

// 新規todoバリデーションの構造体
type NewTodoRequest struct {
	Title         string  `json:"title" validate:"required"`
	Content       string  `json:"content" validate:"required"`
	DueDate       *string `json:"due_date" validate:"omitempty"`
	CompletedDate *string `json:"completed_date" validate:"omitempty"`
}

// 新規todoバリデーション
func (r *NewTodoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
