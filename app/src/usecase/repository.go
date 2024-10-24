package usecase

import (
	"app/src/entities"
	"gorm.io/gorm"
)

// リポジトリ型の構造体を作成⇨レシーバで使用
type Repository struct {
	DB *gorm.DB
}

// DBからのデータ取得やDBへのinsertなど、DB操作を記述する

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
