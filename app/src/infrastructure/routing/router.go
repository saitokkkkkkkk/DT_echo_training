package routing

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// セッションストアを定義
var store = sessions.NewCookieStore([]byte("my-secret-key"))

func Init() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// セッションのミドルウェア
	e.Use(sessionMiddleware)

	// set template
	SetTemplate(e)

	// set routing
	SetRouting(e)

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}

// セッション管理用のミドルウェア
func sessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := store.Get(c.Request(), "session-name")
		if err != nil {
			return err
		}

		// セッションをコンテキストに保存
		c.Set("session", session)

		// 次のハンドラーを呼び出す
		return next(c)
	}
}
