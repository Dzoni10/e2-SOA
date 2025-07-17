package main

import (
	"database-example/database"
	"database-example/handler"
	"database-example/repo"
	"database-example/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.Init()

	repository := &repo.UserRepository{DatabaseConnection: database.DB}
	userService := &service.UserService{Repo: repository}
	userHandler := &handler.UserHandler{UserService: userService}

	router := mux.NewRouter()
	router.HandleFunc("/users/all", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.Get).Methods("GET")
	router.HandleFunc("/users", userHandler.Create).Methods("POST")
	router.HandleFunc("/users/login", userHandler.Login).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("Server starting at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
