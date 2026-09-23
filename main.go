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

	supabaseauth "github.com/supabase-community/auth-go"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Configuration error:", err)
		return
	}

	authClient := supabaseauth.New(
		config.SupabaseProjectRef,
		config.SupabaseKey,
	)

	handlers.AuthClient = authClient

	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}
	handlers.DB = db
	defer db.Close()

	http.HandleFunc("/tasks", handlers.TaskHandler)
	http.HandleFunc("/tasks/", handlers.TaskByIDHandler)

	http.HandleFunc("/auth/signup", handlers.SignupHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)

	http.HandleFunc("/public/info", handlers.PublicInfoHandler)
	http.HandleFunc("/protected/profile", handlers.ProtectedProfileHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Printf("Server running on http://localhost:%s\n", config.Port)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
