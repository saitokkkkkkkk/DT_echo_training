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
	"strings"
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

	// CompletedDateとDueDateを空文字に置換
	if todo.CompletedDate != nil {
		*todo.CompletedDate = strings.Replace(strings.Replace(*todo.CompletedDate, "T", " ", 1), "Z", "", 1)
	}
	if todo.DueDate != nil {
		*todo.DueDate = strings.Replace(strings.Replace(*todo.DueDate, "T", " ", 1), "Z", "", 1)
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

	//バインド前のリクエストの内容をログに出力
	log.Printf("バインド前のrequestの中身: \n %+v", request)

	// ここでrequest変数にフォームからのデータが格納（というかバインド）される
	if err := ctx.Bind(&request); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid input")
	}

	// バインド後のリクエストの内容をログに出力
	log.Printf("バインド後のrequestの中身: \n%+v", request)

	// リクエストの内容をログに出力
	log.Printf("Received request: %+v\n", request)

	// DueDateがnilでない場合、その値をログに出力
	if request.DueDate != nil {
		log.Printf("DueDateの値: %s", *request.DueDate) // ポインタが指している値を表示
	} else {
		log.Println("DueDateはnilです!!")
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

// ToggleTodoStatus handles the toggling of the todo status
func (c *Controller) ToggleTodoStatus(ctx echo.Context) error {
	idParam := ctx.Param("id")
	todoID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Print(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// 完了時間を取得
	completionTime := ctx.FormValue("completionTime")

	// todoの取得
	todo, err := c.Interactor.GetTodoByID(int64(todoID))
	if err != nil {
		log.Print(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Todo not found"})
	}

	// ステータスをトグル
	if completionTime == "" { // undoneに変更
		todo.CompletedDate = nil
	} else { // doneに変更
		todo.CompletedDate = &completionTime
	}

	// データベースに保存
	err = c.Interactor.UpdateTodoStatus(todo)
	if err != nil {
		log.Print(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
	}

	// 新しいステータスを返す
	status := "undone"
	if todo.CompletedDate != nil {
		status = "done"
	}
	return ctx.JSON(http.StatusOK, map[string]string{"status": status})
}

func (c *Controller) GetTodoByID(ctx echo.Context) error {
	idParam := ctx.Param("id")
	todoID, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	todo, err := c.Interactor.GetTodoByID(int64(todoID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	}

	return ctx.JSON(http.StatusOK, todo) // TodoオブジェクトをJSON形式で返す
}

/*会員登録の処理
func (c Controller) RegisterUser(context echo.Context) error {

}*/
