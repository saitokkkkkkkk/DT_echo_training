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
	e.GET("/todos/new", func(c echo.Context) error {
		return c.File("template/new_todo.html")
	})

	// 新規todoを保存
	e.POST("/todos/new", controller.CreateTodo)

	// todo詳細表示
	e.GET("/todos/:id", controller.ShowTodoDetails)

	// 編集画面の表示
	e.GET("/todos/:id/edit", controller.ShowTodoEdit)

	// 編集後の更新
	e.POST("/todos/:id/update", controller.UpdateTodo)

	// todo削除
	e.POST("/todos/:id/delete", controller.DeleteTodo)

	// 一括削除
	e.POST("/todos/bulk-delete", controller.BulkDeleteTodos)

	// done、undoneの更新
	e.PUT("/todos/:id/status", controller.UpdateTodoStatus)
}
