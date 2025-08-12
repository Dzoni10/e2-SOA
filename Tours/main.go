package main

import (
	"log"
	"net/http"
	"tours/database"
	"tours/handler"
	"tours/repo"
	"tours/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	database.Init()

	repo := &repo.TourRepository{}
	srv := &service.TourService{Repo: repo}
	h := &handler.TourHandler{Service: srv}

	r := mux.NewRouter()

	r.HandleFunc("/tours/all", h.GetAllTours).Methods("GET")
	r.HandleFunc("/tours/{id}", h.GetTour).Methods("GET")
	r.HandleFunc("/tours", h.CreateTour).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsHandler.Handler(r)))
}
