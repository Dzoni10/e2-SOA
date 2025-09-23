package main

import (
	"gateway/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CORS middleware (ako ti treba za Angular na :4200)
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func main() {
	router := mux.NewRouter()

	// Handleri
	grpcHandler := &handler.GrpcHandler{}
	gatewayHandler := &handler.GatewayHandler{}

	// --- Blog rute preko gRPC ---
	router.HandleFunc("/blogs/all", grpcHandler.GetAllBlogs).Methods("GET")
	router.HandleFunc("/blogs/{id}", grpcHandler.GetBlog).Methods("GET")
	router.HandleFunc("/blogs", grpcHandler.CreateBlog).Methods("POST")
	router.HandleFunc("/blogs/{id}", grpcHandler.UpdateBlog).Methods("PUT")

	// --- Ostale blog rute (komentari, slike, lajkovi...) ostaju preko proxy-ja ---
	router.PathPrefix("/blogs").HandlerFunc(gatewayHandler.HandleBlog)

	// --- Ostali servisi i dalje preko proxy-ja ---
	router.PathPrefix("/users").HandlerFunc(gatewayHandler.HandleStakeholders)
	router.PathPrefix("/followers").HandlerFunc(gatewayHandler.HandleFollowers)
	router.PathPrefix("/tours").HandlerFunc(gatewayHandler.HandleTours)
	router.PathPrefix("/reviews").HandlerFunc(gatewayHandler.HandleReviews)
	router.PathPrefix("/keypoints").HandlerFunc(gatewayHandler.HandleKeypoints)
	router.PathPrefix("/position").HandlerFunc(gatewayHandler.HandlePosition)
	router.PathPrefix("/tour-executions").HandlerFunc(gatewayHandler.HandleTourExecutions)
	router.PathPrefix("/purchase").HandlerFunc(gatewayHandler.HandlePurchase)

	// Pokreni Gateway
	log.Println("API Gateway started on :8070")
	log.Fatal(http.ListenAndServe(":8070", withCORS(router)))
}
