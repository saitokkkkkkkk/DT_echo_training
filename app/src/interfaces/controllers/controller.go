package controllers

import (
	"app/src/entities"
	"app/src/infrastructure/sqlhandler"
	"app/src/usecase"
	"fmt"
	"github.com/gorilla/sessions"
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

// todoの編集画面の表示（オッケー）
func (c Controller) ShowTodoEdit(ctx echo.Context) error {
	// URLからIDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	// Todo取得
	todo, err := c.Interactor.GetTodoByID(id)
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusNotFound, "Todo not found")
	}

	return ctx.Render(http.StatusOK, "todo_edit.html", todo) // 編集画面のHTMLを返す
}

// todo詳細表示（オッケー）
func (c Controller) ShowTodoDetails(ctx echo.Context) error {
	// URLからIDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	// todo取得
	todo, err := c.Interactor.GetTodoByID(id)
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusNotFound, "Todo not found")
	}

	return ctx.Render(http.StatusOK, "todo_detail.html", todo)
}

// 一括削除（オッケー）
func (c Controller) BulkDeleteTodos(ctx echo.Context) error {
	// Interactorを呼ぶ
	err := c.Interactor.BulkDeleteTodos()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to delete done Todo")
	}

	// セッションを取得
	session := ctx.Get("session").(*sessions.Session)
	// メッセージをセッションに設定
	session.Values["message"] = "All completed todos have been deleted."
	// セッションを保存
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	// 一覧画面にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// 新規todoの保存
func (c Controller) CreateTodo(ctx echo.Context) error {
	// フォームから送信されたデータを取得
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	dueDateStr := ctx.FormValue("due_date")
	completedDateStr := ctx.FormValue("completed_date")

	// Todo構造体を作成
	todo, err := c.prepareTodo(title, content, dueDateStr, completedDateStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// インタラクターを呼び出してTodoを保存
	if err := c.Interactor.CreateTodo(todo); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// セッションを取得
	session := ctx.Get("session").(*sessions.Session)
	session.Values["message"] = "Todo created successfully."
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	// 一覧画面にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// Todoを準備するヘルパー関数
func (c Controller) prepareTodo(title, content, dueDateStr, completedDateStr string) (entities.Todo, error) {
	const layout = "2006-01-02"
	var dueDate, completedDate *time.Time

	// DueDateの処理
	if dueDateStr != "" {
		t, err := time.Parse(layout, dueDateStr)
		if err != nil {
			return entities.Todo{}, fmt.Errorf("Invalid due date format.")
		}
		dueDate = &t
	}

	// CompletedDateの処理
	if completedDateStr != "" {
		t, err := time.Parse(layout, completedDateStr)
		if err != nil {
			return entities.Todo{}, fmt.Errorf("Invalid due date format.")
		}
		completedDate = &t
	}

	// Todo構造体にデータを詰める
	todo := entities.Todo{
		Title:         title,
		Content:       content,
		DueDate:       dueDate,
		CompletedDate: completedDate,
		CreatedAt:     time.Now(), // 作成日時を現在の時刻で設定
	}
	return todo, nil
}

// todoの更新
func (c Controller) UpdateTodo(ctx echo.Context) error {
	// フォームから送信されたデータを取得
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	dueDateStr := ctx.FormValue("due_date")
	completedDateStr := ctx.FormValue("completed_date")

	// URLからtodoのIDを取得して整数に変換
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid ID")
	}

	// 更新するTodoの既存データを取得
	todo, err := c.Interactor.GetTodoByID(id)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to find Todo")
	}

	// 更新後のタイトルと内容をtodoに格納
	todo.Title = title
	todo.Content = content

	// 完了予定日のフォーマット変換
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
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

	// セッションを取得
	session := ctx.Get("session").(*sessions.Session)
	// メッセージをセッションに設定
	session.Values["message"] = "Todo has been successfully updated."
	// セッションを保存
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	// 一覧画面にリダイレクト
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

	// セッションを取得
	session := ctx.Get("session").(*sessions.Session)
	// メッセージをセッションに設定
	session.Values["message"] = "Todo has been successfully deleted."
	// セッションを保存
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	// 一覧画面にリダイレクト
	return ctx.Redirect(http.StatusSeeOther, "/todos")
}

// doneとundoneのステータス更新 ⇨ 一覧画面表示
func (c Controller) UpdateTodoStatus(ctx echo.Context) error {

	// 変更されたステータスをリクエストボディ（json）から取得して構造体にバインド
	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := ctx.Bind(&statusUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"ok":    false, //jsのokのとこに返る
			"error": "Bind Error",
		})
	}

	// URLからidを取得してintに
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"ok":    false,
			"error": "Invalid ID",
		})
	}

	// インタラクターを使用してステータス更新
	err = c.Interactor.UpdateTodoStatus(id, statusUpdate.Status)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"ok":    false,
			"error": "Failed to update status",
		})
	}

	// ステータス更新後に全てのTodoを取得
	todos, err := c.Interactor.GetAllTodos()
	if err != nil {
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
		"hasCompleted": hasCompleted, // 完了したTodoがあるかどうかを返却
	})
}

// todo一覧表示
func (c Controller) Index(ctx echo.Context) error {
	// todoの取得
	todos, err := c.Interactor.GetAllTodos()
	if err != nil {
		log.Print(err)
		return ctx.Render(http.StatusInternalServerError, "todo_list.html", nil)
	}

	// セッションの取得
	session := ctx.Get("session").(*sessions.Session)
	// メッセージを取得
	message := session.Values["message"]
	// メッセージを削除（=リクエストごとにメッセージがあれば設定）
	delete(session.Values, "message")
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	// 完了したTodoが1つ以上あるか確認
	hasCompleted := false
	for _, todo := range todos {
		if todo.CompletedDate != nil {
			hasCompleted = true
			break
		}
	}

	// マップに変換してデータを画面に渡す
	return ctx.Render(http.StatusOK, "todo_list.html", map[string]interface{}{
		"Todos":        todos,
		"HasCompleted": hasCompleted, // 完了済みTodoが1つ以上あるかを渡す
		"Message":      message,      // メッセージを渡す
	})
}
