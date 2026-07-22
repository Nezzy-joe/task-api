package handlers

import (
	"ai-task-api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var tasks = []models.Task{
	{
		ID:        1,
		Title:     "Learn Go",
		Completed: false,
	},
	{
		ID:        2,
		Title:     "Build Task API",
		Completed: false,
	},
}

func getTaskID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/tasks/")
	return strconv.Atoi(idStr)
}

// HandleTasks handles GET /tasks and POST /tasks
func HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask models.Task

		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(newTask.Title) == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		newTask.ID = len(tasks) + 1
		tasks = append(tasks, newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleTaskByID handles GET, PUT and DELETE for a single task.
func HandleTaskByID(w http.ResponseWriter, r *http.Request) {

	id, err := getTaskID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {

		if task.ID == id {

			switch r.Method {

			case http.MethodGet:

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(task)
				return

			case http.MethodPut:

				var updatedTask models.Task

				if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				if strings.TrimSpace(updatedTask.Title) == "" {
					http.Error(w, "Title is required", http.StatusBadRequest)
					return
				}

				updatedTask.ID = id
				tasks[i] = updatedTask

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedTask)
				return

			case http.MethodDelete:

				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return

			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
