package main

import (
	"context"
	"follower-service/database"
	"follower-service/handler"
	"follower-service/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	if err := database.InitNeo4j(); err != nil {
		log.Fatal("Failed to initialize Neo4j:", err)
	}

	followerService := service.NewFollowerService()
	followerHandler := handler.NewFollowerHandler(followerService)

	r := mux.NewRouter()

	r.HandleFunc("/users", followerHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{userId}/followers", followerHandler.GetFollowers).Methods("GET")
	r.HandleFunc("/users/{userId}/following", followerHandler.GetFollowing).Methods("GET")
	r.HandleFunc("/users/{userId}/stats", followerHandler.GetUserStats).Methods("GET")
	r.HandleFunc("/users/{userId}/recommendations", followerHandler.GetFollowRecommendations).Methods("GET")

	// Follow/Unfollow routes
	r.HandleFunc("/follow", followerHandler.FollowUser).Methods("POST")
	r.HandleFunc("/unfollow", followerHandler.UnfollowUser).Methods("POST")

	// Check relationships
	r.HandleFunc("/is-following/{followerId}/{followingId}", followerHandler.IsFollowing).Methods("GET")
	r.HandleFunc("/can-comment/{commenterId}/{authorId}", followerHandler.CanUserComment).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Follower Service is running"))
	}).Methods("GET")

	// CORS setup
	/*corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:8080", "http://localhost:8082"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})*/

	// Get port from environment or default to 8083
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Follower Service starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close Neo4j connection
	database.CloseNeo4j(ctx)
	log.Println("Server exited")
}
