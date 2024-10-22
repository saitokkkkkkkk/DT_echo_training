package controllers

import (
	"app/src/entities"
	"app/src/infrastructure/sqlhandler"
	"app/src/usecase"
	"app/src/validation"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

// コントローラ型の構造体を作成⇨レシーバとして使用
type Controller struct {
	Interactor usecase.Interactor
}

// コントローラのインスタンスを作成
func NewController(sqlhandler *sqlhandler.SqlHandler) *Controller {
	return &Controller{
		Interactor: usecase.Interactor{
			Repository: usecase.Repository{
				DB: sqlhandler.DB,
			},
		},
	}
}

// todo一覧表示
func (c Controller) Index(ctx echo.Context) error {
	todos, err := c.Interactor.GetAllTodos()
	if err != nil {
		log.Print(err)
		return ctx.Render(http.StatusInternalServerError, "todo_list.html", nil)
	}

	// todosをマップに変換して渡す
	return ctx.Render(http.StatusOK, "todo_list.html", map[string]interface{}{
		"Todos": todos, // "Todos"というキーで渡す
	})
}

// todo詳細表示
func (c Controller) ShowTodoDetails(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id")) // URLからIDを取得
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	todo, err := c.Interactor.GetTodoByID(int64(id)) // Interactorを使ってTodoを取得
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusNotFound, "Todo not found")
	}

	return ctx.Render(http.StatusOK, "todo_detail.html", todo) // 詳細表示用のHTMLを返す
}

// todo新規作成画面の表示
func (c Controller) ShowNewTodoForm(ctx echo.Context) error {
	showCompletedDate := false
	return ctx.Render(http.StatusOK, "new_todo.html", map[string]interface{}{
		"ShowCompletedDate": showCompletedDate,
	})
}

// 新規todoの保存
func (c Controller) CreateTodo(ctx echo.Context) error {
	// フォームデータに対してバリデーション
	var request validators.NewTodoRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid input")
	}

	// リクエストの内容をログに出力
	log.Printf("Received request: %+v\n", request)

	// DueDateがnilでない場合、その値をログに出力
	if request.DueDate != nil {
		log.Printf("DueDate value: %s", *request.DueDate) // ポインタが指している値を表示
	} else {
		log.Println("DueDate is nil")
	}

	// バリデーション
	if err := request.Validate(); err != nil {
		return ctx.Render(http.StatusBadRequest, "new_todo.html", map[string]interface{}{
			"ShowCompletedDate": false,
			"ValidationErrors":  []string{err.Error()},
		})
	}

	// 保存するnewTodoを作成
	newTodo := entities.Todo{
		Title:   request.Title,
		Content: request.Content,
	}

	//DueDateがnilでも空文字でもなければそのまま代入
	if *request.DueDate == "" {
		newTodo.DueDate = nil
	} else {
		newTodo.DueDate = request.DueDate
	}

	// Interactorを使用して新規Todoを保存
	if err := c.Interactor.CreateTodo(newTodo); err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to create Todo")
	}

	// 保存成功後は一覧にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}
