package main

import (
	"fmt"
	"net/http"

	"github.com/Nezzy-joe/task-api/handlers"
)

func main() {
	http.HandleFunc("/tasks", handlers.TaskHandler)
	http.HandleFunc("/tasks/", handlers.GetTaskByID)
	fmt.Println("server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
