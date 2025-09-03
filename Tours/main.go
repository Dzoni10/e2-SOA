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

	tourRepo := &repo.TourRepository{}
	srv := &service.TourService{Repo: tourRepo}
	h := &handler.TourHandler{Service: srv}

	reviewRepo := &repo.ReviewRepository{}
	reviewSrv := &service.ReviewService{Repo: reviewRepo}
	reviewHandler := &handler.ReviewHandler{Service: reviewSrv}

	r := mux.NewRouter()

	r.HandleFunc("/tours/all", h.GetAllTours).Methods("GET")
	r.HandleFunc("/tours/{id}", h.GetTour).Methods("GET")
	r.HandleFunc("/tours", h.CreateTour).Methods("POST")

	r.HandleFunc("/tours/{id}/reviews", reviewHandler.GetReviewsForTour).Methods("GET")
	r.HandleFunc("/reviews", reviewHandler.AddReview).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsHandler.Handler(r)))
}
