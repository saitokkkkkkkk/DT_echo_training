package usecase

import (
	"app/entities"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// このファイルではDBからのデータ取得やDBへのinsertなど、DB操作を記述する

func(r *Repository) GetAllArticle() (articles []entities.Article, err error) {
	// 以下は実際にはDBを使って記事の全データを取得したりする
	var article entities.Article
	article.ID = 1
	article.Title = "Deep Track"
	articles = append(articles, article)
	return articles, nil
}
