package usecase

import (
	"app/src/entities"
)

type Interactor struct {
	Repository Repository
}

// アプリケーション固有のビジネスルール
// このファイルでは取得したデータを組み合わせたりしてユースケースを実現する

func(i Interactor) GetAllArticle() (article []entities.Article, err error) {
	return i.Repository.GetAllArticle()
}