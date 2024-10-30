package usecase

import (
	"app/src/entities"
	"errors"
	"time"
)

type Interactor struct {
	Repository Repository
}

// todo一覧取得
func (i Interactor) GetAllTodos() (todos []entities.Todo, err error) {

	return i.Repository.GetAllTodos()
}

// todo取得
func (i Interactor) GetTodoByID(id int) (entities.Todo, error) {

	return i.Repository.GetTodoByID(id)
}

// 新規todo保存
func (i Interactor) CreateTodo(todo entities.Todo) error {

	return i.Repository.CreateTodo(todo)
}

// Todoのステータスを更新
func (i *Interactor) UpdateTodo(todo entities.Todo) error {
	return i.Repository.UpdateTodo(todo)
}

// todo単体削除
func (i *Interactor) DeleteTodo(id int64) error {
	// Repository の DeleteTodo メソッドを呼び出し、エラーをそのまま返す
	return i.Repository.DeleteTodo(id)
}

// todo一括削除
func (i *Interactor) BulkDeleteTodos() error {
	// 条件をインタラクタで定義（リポジトリに書いた方が良いのか）
	condition := "completed_date IS NOT NULL"
	return i.Repository.BulkDeleteTodos(condition)
}

// done, undoneのステータス更新
func (i *Interactor) UpdateTodoStatus(id int, status string) error {
	if status == "done" {
		return i.Repository.SetCompletedAt(id, time.Now())
	} else if status == "undone" {
		return i.Repository.SetCompletedAtNull(id)
	}
	return errors.New("Invalid status or id")
}
