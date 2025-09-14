package main

import (
	"gateway/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Ručni CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obriši sve eventualne duple vrednosti
		w.Header().Del("Access-Control-Allow-Origin")

		// Postavi samo jedan
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func startServer(handler *handler.GatewayHandler) {
	router := mux.NewRouter()

	// tvoje rute
	router.PathPrefix("/blogs").HandlerFunc(handler.HandleBlog)
	router.PathPrefix("/users").HandlerFunc(handler.HandleStakeholders)
	router.PathPrefix("/followers").HandlerFunc(handler.HandleFollowers)
	router.PathPrefix("/tours").HandlerFunc(handler.HandleTours)
	router.PathPrefix("/reviews").HandlerFunc(handler.HandleReviews)
	router.PathPrefix("/keypoints").HandlerFunc(handler.HandleKeypoints)
	router.PathPrefix("/position").HandlerFunc(handler.HandlePosition)

	log.Println("Gateway started on :8070")
	log.Fatal(http.ListenAndServe(":8070", withCORS(router)))
}

func main() {
	h := &handler.GatewayHandler{}
	startServer(h)
}
