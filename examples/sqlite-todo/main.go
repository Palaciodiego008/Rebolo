package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo"
	"github.com/gorilla/mux"
)

// Todo model
type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var app *rebolo.Application

func main() {
	// Create new ReboloLang app
	// Database config is loaded from config.yml
	app = rebolo.New()

	// Run migrations
	if err := migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Define routes
	app.GET("/", homeHandler)
	app.GET("/todos", listTodosHandler)
	app.POST("/todos", createTodoHandler)
	app.GET("/todos/{id}", getTodoHandler)
	app.PUT("/todos/{id}", updateTodoHandler)
	app.DELETE("/todos/{id}", deleteTodoHandler)

	// Start server
	log.Println("🚀 Starting Todo API with SQLite...")
	if err := app.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// migrate creates the database tables
func migrate() error {
	db := app.DB()
	ctx := context.Background()

	// Create todos table using standard SQL
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrated successfully")
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	app.RenderJSON(w, map[string]string{
		"message": "Welcome to ReboloLang Todo API with SQLite!",
		"version": "1.0",
	})
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	db := app.DB()

	rows, err := db.QueryContext(r.Context(),
		"SELECT id, title, completed, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		app.RenderError(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt); err != nil {
			app.RenderError(w, "Failed to scan todos", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	app.RenderJSON(w, todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	db := app.DB()

	if err := r.ParseForm(); err != nil {
		app.RenderError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		app.RenderError(w, "Title is required", http.StatusBadRequest)
		return
	}

	result, err := db.ExecContext(r.Context(),
		"INSERT INTO todos (title, completed, created_at) VALUES (?, ?, ?)",
		title, false, time.Now())
	if err != nil {
		app.RenderError(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	todo := Todo{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	app.RenderJSON(w, todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	db := app.DB()
	vars := mux.Vars(r)
	id := vars["id"]

	var todo Todo
	err := db.QueryRowContext(r.Context(),
		"SELECT id, title, completed, created_at FROM todos WHERE id = ?", id).
		Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)

	if err != nil {
		app.RenderError(w, "Todo not found", http.StatusNotFound)
		return
	}

	app.RenderJSON(w, todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	db := app.DB()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := r.ParseForm(); err != nil {
		app.RenderError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// First, get the existing todo
	var todo Todo
	err := db.QueryRowContext(r.Context(),
		"SELECT id, title, completed, created_at FROM todos WHERE id = ?", id).
		Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)

	if err != nil {
		app.RenderError(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if title := r.FormValue("title"); title != "" {
		todo.Title = title
	}
	if completed := r.FormValue("completed"); completed != "" {
		todo.Completed = completed == "true"
	}

	// Update in database
	_, err = db.ExecContext(r.Context(),
		"UPDATE todos SET title = ?, completed = ? WHERE id = ?",
		todo.Title, todo.Completed, id)

	if err != nil {
		app.RenderError(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	app.RenderJSON(w, todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	db := app.DB()
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.ExecContext(r.Context(), "DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		app.RenderError(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}
