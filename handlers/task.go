package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nezzy-joe/task-api/models"
)

//memory storage for tasks.

var Tasks = []models.Task{
	{
		ID:        1,
		Title:     "Learn Go",
		Completed: false,
	},
	{
		ID:        2,
		Title:     "Build Task Api",
		Completed: false,
	},
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Tasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, task := range Tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	newTask.ID = len(Tasks) + 1

	Tasks = append(Tasks, newTask)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTask(w, r)

	case http.MethodPost:
		CreateTask(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
