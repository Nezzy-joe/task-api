package main

import (
	"fmt"
	"net/http"

	"ai-task-api/handlers"
)

func main() {
	http.HandleFunc("/tasks", handlers.HandleTasks)
	http.HandleFunc("/tasks/", handlers.HandleTaskByID)

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
