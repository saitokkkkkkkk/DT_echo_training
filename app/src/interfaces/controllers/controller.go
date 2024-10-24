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
	"time"
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

	// Interactorを使用して新規Todoを保存
	if err := c.Interactor.CreateTodo(newTodo); err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to create Todo")
	}

	// 保存成功後は一覧にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// todoの編集画面の表示
func (c Controller) ShowTodoEdit(ctx echo.Context) error {
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

	return ctx.Render(http.StatusOK, "todo_edit.html", todo) // 編集画面のHTMLを返す
}

// todoの更新
func (c Controller) UpdateTodo(ctx echo.Context) error {
	// URLからIDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	// 更新するTodoの既存データを取得
	todo, err := c.Interactor.GetTodoByID(int64(id))
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to find Todo")
	}

	// フォームから送信されたデータを取得
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	dueDateStr := ctx.FormValue("due_date")
	completedDateStr := ctx.FormValue("completed_date")

	// フォームデータをTodoに反映
	todo.Title = title
	todo.Content = content

	// 日付データの変換
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02T15:04", dueDateStr)
		if err != nil {
			log.Print(err)
			return ctx.String(http.StatusBadRequest, "Invalid Due Date format")
		}
		todo.DueDate = &dueDate
	} else {
		todo.DueDate = nil // 日付が未入力の場合、nilにする
	}

	if completedDateStr != "" {
		completedDate, err := time.Parse("2006-01-02T15:04", completedDateStr)
		if err != nil {
			log.Print(err)
			return ctx.String(http.StatusBadRequest, "Invalid Completed Date format")
		}
		todo.CompletedDate = &completedDate
	} else {
		todo.CompletedDate = nil // 日付が未入力の場合、nilにする
	}

	// Interactorを使用してTodoを更新
	if err := c.Interactor.UpdateTodo(todo); err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to update Todo")
	}

	// 更新後、一覧ページにリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// todo削除
func (c Controller) DeleteTodo(ctx echo.Context) error {
	// URLからIDを取得
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid Todo ID")
	}

	// Interactorを使用して削除処理を実行
	err = c.Interactor.DeleteTodo(int64(id))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete Todo")
	}

	/* セッションを取得
	session := ctx.Get("session").(*sessions.Session)
	// 成功メッセージをセッションに設定
	session.Values["message"] = "Deleted successfully"
	session.Save(ctx.Request(), ctx.Response())*/

	// 成功時、Todo一覧にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// done todo 一括削除
func (c Controller) BulkDeleteTodos(ctx echo.Context) error {
	// Interactorを呼ぶ
	err := c.Interactor.BulkDeleteTodos()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete done Todo")
	}

	// 一括削除成功後、Todo一覧にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

/*会員登録の処理
func (c Controller) RegisterUser(context echo.Context) error {

}*/
