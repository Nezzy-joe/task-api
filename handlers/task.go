package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nezzy-joe/task-api/models"
)

var DB *sql.DB

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// GetTasks godoc
//
//	@Summary		Get all tasks
//	@Description	Return all tasks
//	@Tags			Tasks
//	@Produce		json
//	@Success		200	{array}	models.Task
//	@Router			/tasks [get]
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := DB.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		writeJSONError(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task
		var done int

		if err := rows.Scan(&task.ID, &task.Title, &done); err != nil {
			writeJSONError(w, "Failed to read tasks", http.StatusInternalServerError)
			return
		}

		task.Completed = done == 1
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		writeJSONError(w, "Failed to read tasks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID godoc
//
//	@Summary		Get task by ID
//	@Description	Return a single task
//	@Tags			Tasks
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	models.Task
//	@Failure		400	{string}	string	"Invalid task ID"
//	@Failure		404	{string}	string	"Task not found"
//	@Router			/tasks/{id} [get]
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	var done int

	err = DB.QueryRow(
		"SELECT id, title, done FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Title, &done)

	if err == sql.ErrNoRows {
		writeJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	task.Completed = done == 1

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask godoc
//
//	@Summary		Update a task
//	@Description	Update an existing task
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Task ID"
//	@Param			task	body		models.Task	true	"Updated Task"
//	@Success		200		{object}	models.Task
//	@Failure		400		{string}	string	"Invalid JSON"
//	@Failure		404		{string}	string	"Task not found"
//	@Router			/tasks/{id} [put]
func UpdateTask(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return

	}

	var UpdatedTask models.Task
	err = json.NewDecoder(r.Body).Decode(&UpdatedTask)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if UpdatedTask.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	done := 0
	if UpdatedTask.Completed {
		done = 1
	}

	result, err := DB.Exec(
		"UPDATE tasks SET title = ?, done = ? WHERE id = ?",
		UpdatedTask.Title,
		done,
		id,
	)

	if err != nil {
		writeJSONError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Failed to verify update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	UpdatedTask.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UpdatedTask)
}

// DeleteTask godoc
//
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Tags			Tasks
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204
//	@Failure		400	{string}	string	"Invalid task ID"
//	@Failure		404	{string}	string	"Task not found"
//	@Router			/tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result, err := DB.Exec(
		"DELETE FROM tasks WHERE id = ?",
		id,
	)

	if err != nil {
		writeJSONError(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Failed to verify deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateTask godoc
//
//	@Summary		Create a new task
//	@Description	Create a new task
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.Task	true	"Task"
//	@Success		201		{object}	models.Task
//	@Failure		400		{string}	string	"Invalid JSON"
//	@Router			/tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if newTask.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}
	done := 0

	if newTask.Completed {
		done = 1
	}

	result, err := DB.Exec(
		"INSERT INTO tasks (title, done) VALUES (?, ?)",
		newTask.Title,
		done,
	)

	if err != nil {
		writeJSONError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		writeJSONError(w, "Failed to get created task ID", http.StatusInternalServerError)
		return
	}

	newTask.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// TASK HANDLER
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTask(w, r)

	case http.MethodPost:
		CreateTask(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//TASKBYIDHANDLER

func TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTaskByID(w, r)

	case http.MethodPut:
		UpdateTask(w, r)

	case http.MethodDelete:
		DeleteTask(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
