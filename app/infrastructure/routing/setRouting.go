package routing

import (
	"app/infrastructure/sqlhandler"
	"app/interfaces/controllers"
	"github.com/labstack/echo/v4"
	"net/http"
)

// このファイルにはリクエストのルーティング処理を実装する

func SetRouting(e *echo.Echo) {

	controller := controllers.NewController(sqlhandler.NewSqlHandler())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo World!!")
	})

	e.GET("/allArticles", func(c echo.Context) error {
		return controller.Index(c)
	})
}
