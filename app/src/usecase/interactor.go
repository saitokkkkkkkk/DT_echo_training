package usecase

import (
	"app/src/entities"
)

type Interactor struct {
	Repository Repository
}

// アプリケーション固有のビジネスルール
// このファイルでは取得したデータを組み合わせたりしてユースケースを実現する

// todo一覧取得
func (i Interactor) GetAllTodos() (todos []entities.Todo, err error) {

	return i.Repository.GetAllTodos()
}

// todo取得
func (i Interactor) GetTodoByID(id int64) (entities.Todo, error) {

	return i.Repository.GetTodoByID(id)
}

// 新規todo保存
func (i Interactor) CreateTodo(todo entities.Todo) error {
	// Todoを保存するためのロジックをここに書く
	return i.Repository.CreateTodo(todo)
}

// todoのステータスを更新
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
	return i.Repository.BulkDeleteTodos()
}
