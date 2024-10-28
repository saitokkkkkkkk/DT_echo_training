package controllers

import (
	"app/src/entities"
	"app/src/infrastructure/sqlhandler"
	"app/src/usecase"
	"fmt"
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

	// 完了したTodoが1つ以上あるか確認
	hasCompleted := false
	for _, todo := range todos {
		if todo.CompletedDate != nil {
			hasCompleted = true
			break
		}
	}

	// todosをマップに変換して渡す
	return ctx.Render(http.StatusOK, "todo_list.html", map[string]interface{}{
		"Todos":        todos,
		"HasCompleted": hasCompleted, // 完了済みTodoが1つ以上あるかを渡す
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

// 新規todoの保存
func (c Controller) CreateTodo(ctx echo.Context) error {

	// フォームから送信されたデータを取得
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	dueDateStr := ctx.FormValue("due_date")
	completedDateStr := ctx.FormValue("completed_date")

	// 日付のフォーマットを指定
	const layout = "2006-01-02"
	var dueDate, completedDate *time.Time

	// DueDateの処理
	if dueDateStr != "" {
		t, err := time.Parse(layout, dueDateStr)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "無効な期限日フォーマットです")
		}
		dueDate = &t
	}

	// CompletedDateの処理
	if completedDateStr != "" {
		t, err := time.Parse(layout, completedDateStr)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "無効な完了日フォーマットです")
		}
		completedDate = &t
	}

	// Todo構造体にフォームデータを詰める
	todo := entities.Todo{
		Title:         title,
		Content:       content,
		DueDate:       dueDate,
		CompletedDate: completedDate,
		CreatedAt:     time.Now(), // 作成日時を現在の時刻で設定
	}

	// 詰めたtodoをインタラクターに渡して呼ぶ
	if err := c.Interactor.CreateTodo(todo); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 成功した場合は/todosにリダイレクト
	return ctx.Redirect(http.StatusFound, "/todos")
}

// todoの更新
func (c Controller) UpdateTodo(ctx echo.Context) error {

	// フォームから送信されたデータを取得
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	dueDateStr := ctx.FormValue("due_date")
	completedDateStr := ctx.FormValue("completed_date")

	// URLからtodoのIDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	// 更新するTodoの既存データを取得
	todo, err := c.Interactor.GetTodoByID(int64(id))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to find Todo")
	}

	// タイトルと内容を更新
	todo.Title = title
	todo.Content = content

	// 完了予定日のフォーマット変換
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			log.Print(err)
			return ctx.String(http.StatusBadRequest, "Invalid Due Date format")
		}
		todo.DueDate = &dueDate
		log.Printf("パース後のデータ（DueDate）: %v (type=%T)", todo.DueDate, todo.DueDate)
	} else {
		todo.DueDate = nil // 日付が未入力の場合、nilにする
		log.Println("Due date is empty, setting DueDate to nil")
	}

	// 完了日のフォーマット変換
	if completedDateStr != "" {
		completedDate, err := time.Parse("2006-01-02", completedDateStr)
		if err != nil {
			log.Print(err)
			return ctx.String(http.StatusBadRequest, "Invalid Completed Date format")
		}
		todo.CompletedDate = &completedDate
		log.Printf("パース後のデータ（CompletedDate）: %v (type=%T)", todo.CompletedDate, todo.CompletedDate)
	} else {
		todo.CompletedDate = nil // 日付が未入力の場合、nilにする
		log.Println("Completed date is empty, setting CompletedDate to nil")
	}

	// Interactorを使用してTodoを更新
	if err := c.Interactor.UpdateTodo(todo); err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to update Todo")
	}

	// 更新後、一覧ページにリダイレクト
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

// doneとundoneのステータス更新
func (c Controller) UpdateTodoStatus(ctx echo.Context) error {
	// URLからidを取得
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID:", idStr)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"ok":    false,
			"error": "Invalid ID",
		})
	}

	// int⇨int64にする
	todoId := int64(id)

	// 変更されたステータスをリクエストボディから取得して構造体にバインド
	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := ctx.Bind(&statusUpdate); err != nil {
		fmt.Println("Bind error:", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"ok":    false,
			"error": "Invalid ID",
		})
	}

	// インタラクターを使用してステータス更新
	err = c.Interactor.UpdateTodoStatus(todoId, statusUpdate.Status)
	if err != nil {
		fmt.Println("Update status error:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"ok":    false,
			"error": "Failed to update status",
		})
	}

	// ステータス更新後に全てのTodoを取得
	todos, err := c.Interactor.GetAllTodos()
	if err != nil {
		log.Print(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"ok":    false,
			"error": "Failed to fetch todos",
		})
	}

	// 完了したTodoが1つ以上あるか確認
	hasCompleted := false
	for _, todo := range todos {
		if todo.CompletedDate != nil {
			hasCompleted = true
			break
		}
	}

	// 更新結果とhasCompletedを返却
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"ok":           true,
		"status":       statusUpdate.Status,
		"hasCompleted": hasCompleted, // 完了したTodoの状態を返す
	})
}

/*会員登録の処理
func (c Controller) RegisterUser(context echo.Context) error {

}*/
