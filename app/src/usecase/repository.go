package usecase

import (
	"app/src/entities"
	"gorm.io/gorm"
	"time"
)

// リポジトリ型の構造体を作成⇨レシーバで使用
type Repository struct {
	DB *gorm.DB
}

// todo一覧取得
func (r *Repository) GetAllTodos() (todos []entities.Todo, err error) {
	// CreatedAtで昇順にソートして全てのTodoを取得
	if err := r.DB.Order("created_at asc").Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

// todo取得
func (r *Repository) GetTodoByID(id int) (entities.Todo, error) {
	var todo entities.Todo
	if err := r.DB.First(&todo, id).Error; err != nil {
		return todo, err
	}
	return todo, nil
}

// 新規Todo保存
func (r *Repository) CreateTodo(todo entities.Todo) error {
	if err := r.DB.Create(&todo).Error; err != nil {
		return err
	}
	return nil
}

// Todoのステータスを更新
func (r *Repository) UpdateTodo(todo entities.Todo) error {
	if err := r.DB.Save(&todo).Error; err != nil {
		return err
	}
	return nil
}

// 削除
func (r *Repository) DeleteTodo(id int64) error {
	if err := r.DB.Delete(&entities.Todo{}, id).Error; err != nil {
		return err // エラーがあれば返す
	}
	return nil // 正常終了
}

// 一括削除
func (r *Repository) BulkDeleteTodos(condition string) error {
	// 条件に基づいてレコードを削除
	result := r.DB.Where(condition).Delete(&entities.Todo{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// completed_atを設定するメソッド
func (r *Repository) SetCompletedAt(id int, completedAt time.Time) error {
	todo := &entities.Todo{
		ID:            id,
		CompletedDate: &completedAt, // ポインタを使用して、NULLが必要な場合の対応
	}
	if err := r.DB.Model(todo).Update("completed_date", todo.CompletedDate).Error; err != nil {
		return err // エラーがあれば返す
	}
	return nil // 正常終了
}

// completed_atをNULLに設定するメソッド
func (r *Repository) SetCompletedAtNull(id int) error {
	todo := &entities.Todo{ID: id}
	if err := r.DB.Model(todo).Update("completed_date", nil).Error; err != nil {
		return err // エラーがあれば返す
	}
	return nil // 正常終了
}
