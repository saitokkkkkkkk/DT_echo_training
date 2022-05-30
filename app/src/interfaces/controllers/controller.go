package controllers

import (
	"app/src/infrastructure/sqlhandler"
	"app/src/usecase"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)


type Controller struct {
	Interactor usecase.Interactor
}

/*
このファイルには外部からのリクエストで受け取ったデータをusecaseで使えるように変形したり、
内部からのデータを外部機能に向けて便利な形式に変換したりする
例)　外部からのデータをArticleエンティティに変換
*/

func NewController(sqlhandler *sqlhandler.SqlHandler) *Controller {
	return &Controller{
		Interactor: usecase.Interactor{
			Repository: usecase.Repository{
				DB: sqlhandler.DB,
			},
		},
	}
}

func (c Controller) Index(ctx echo.Context) error {
	articles, err := c.Interactor.GetAllArticle()
	if err != nil {
		log.Print(err)
		return ctx.Render(500, "article_list.html", nil)
	}
	return ctx.Render(http.StatusOK, "article_list.html", articles)
}
