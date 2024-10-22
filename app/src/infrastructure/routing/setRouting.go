package routing

import (
	"app/src/infrastructure/sqlhandler"
	"app/src/interfaces/controllers"
	"github.com/labstack/echo/v4"
)

// このファイルにはリクエストのルーティング処理を実装する

func SetRouting(e *echo.Echo) {
	controller := controllers.NewController(sqlhandler.NewSqlHandler())

	// todo一覧表示
	e.GET("/todos", controller.Index)

	// todo新規作成画面を表示
	e.GET("/todos/new", controller.ShowNewTodoForm)

	// 新規todoを保存
	e.GET("/todos/new", controller.CreateTodo)

	// todo詳細表示
	e.GET("/todos/:id", controller.ShowTodoDetails)
}
