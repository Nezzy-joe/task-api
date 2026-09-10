// @title Task API
// @version 1.0
// @description A simple Task Management API built with Go.
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"net/http"

	"github.com/Nezzy-joe/task-api/database"
	"github.com/Nezzy-joe/task-api/handlers"

	_ "github.com/Nezzy-joe/task-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}
	handlers.DB = db
	defer db.Close()

	http.HandleFunc("/tasks", handlers.TaskHandler)
	http.HandleFunc("/tasks/", handlers.TaskByIDHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
