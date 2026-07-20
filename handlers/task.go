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

// GET TASK
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Tasks)
}

// GET TASK BYID
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

// UPDATE TASK
func UpdateTask(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return

	}

	var UpdatedTask models.Task
	err = json.NewDecoder(r.Body).Decode(&UpdatedTask)
	if err != nil {
		http.Error(w, "Invalide JSON", http.StatusBadRequest)
		return
	}
	if UpdatedTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	for i, task := range Tasks {
		if task.ID == id {
			UpdatedTask.ID = id
			Tasks[i] = UpdatedTask
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(UpdatedTask)
			return
		}

	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// DELETE TASK
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	for i, task := range Tasks {
		if task.ID == id {
			Tasks = append(Tasks[:1], Tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}

	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// CREATE TASK
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalide JSON", http.StatusBadRequest)
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

// TASK HANDLER
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
