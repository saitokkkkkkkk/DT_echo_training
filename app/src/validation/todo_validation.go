package validators

// 新規todoバリデーションの構造体
type NewTodoRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	DueDate string `json:"due_date" validate:"omitempty"` // 省略可能
}

// 新規todoバリデーション
func (r *NewTodoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// 更新todoバリデーションの構造体（例）
type UpdateTodoRequest struct {
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	DueDate       string `json:"due_date" validate:"omitempty"`
	CompletedDate string `json:"completed_date" validate:"omitempty"`
}

// 更新todoバリデーション
func (r *UpdateTodoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
