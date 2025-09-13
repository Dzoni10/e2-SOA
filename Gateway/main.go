package main

import (
	"gateway/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func startServer(handler *handler.GatewayHandler) {
	router := mux.NewRouter()

	// tvoje rute
	router.PathPrefix("/blogs").HandlerFunc(handler.HandleBlog)
	router.PathPrefix("/users").HandlerFunc(handler.HandleStakeholders)
	router.PathPrefix("/followers").HandlerFunc(handler.HandleFollowers)
	router.PathPrefix("/tours").HandlerFunc(handler.HandleTours)
	router.PathPrefix("/reviews").HandlerFunc(handler.HandleTours)
	router.PathPrefix("/keypoints").HandlerFunc(handler.HandleTours)
	router.PathPrefix("/position").HandlerFunc(handler.HandleTours)

	// ✅ CORS ovde
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Angular app
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	log.Println("Gateway started on :8070")
	log.Fatal(http.ListenAndServe(":8070", corsHandler.Handler(router)))
}

func main() {
	h := &handler.GatewayHandler{}
	startServer(h)
}
