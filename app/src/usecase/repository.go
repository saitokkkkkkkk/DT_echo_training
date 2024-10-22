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

// todo詳細取得
func (r *Repository) GetTodoByID(id int64) (entities.Todo, error) {
	var todo entities.Todo
	if err := r.DB.First(&todo, id).Error; err != nil {
		return todo, err
	}
	return todo, nil
}
