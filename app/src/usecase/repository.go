package usecase

import (
	"app/src/entities"
	"gorm.io/gorm"
)

// リポジトリ型の構造体を作成⇨レシーバで使用
type Repository struct {
	DB *gorm.DB
}

// todo一覧取得
func (r *Repository) GetAllTodos() (todos []entities.Todo, err error) {
	if err := r.DB.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

// todo取得
func (r *Repository) GetTodoByID(id int64) (entities.Todo, error) {
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

// todoのステータスを更新
func (r *Repository) UpdateTodo(todo entities.Todo) error {
	if err := r.DB.Save(&todo).Error; err != nil {
		return err
	}
	return nil
}

// 一括削除
func (r *Repository) BulkDeleteTodos() error {
	// completed_dateがNULLでないレコードを削除
	result := r.DB.Where("completed_date IS NOT NULL").Delete(&entities.Todo{})
	if result.Error != nil {
		return result.Error // エラーが発生した場合は返す
	}
	return nil // 成功した場合はnilを返す
}

// 削除
func (r *Repository) DeleteTodo(id int64) error {
	if err := r.DB.Delete(&entities.Todo{}, id).Error; err != nil {
		return err // エラーがあれば返す
	}
	return nil // 正常終了
}
