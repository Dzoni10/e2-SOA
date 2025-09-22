package main

import (
	"context"
	"encoding/json"
	"fmt"
	"follower-service/database"
	"follower-service/handler"
	"follower-service/notifications"
	"follower-service/service"
	"local/common/saga/events"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	notifications.Register(userID, conn)
	log.Println("User connected via WS:", userID)

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	notifications.Unregister(userID)
	conn.Close()
	log.Println("User disconnected WS:", userID)
}

func main() {
	if err := database.InitNeo4j(); err != nil {
		log.Fatal("Failed to initialize Neo4j:", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://nats:4222" // default za Docker
	}
	nc, err := events.ConnectNATS(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	followerService := service.NewFollowerService()
	followerHandler := handler.NewFollowerHandler(followerService)

	// Subscribe to blog.created in a goroutine
	go func() {
		_, err := nc.Subscribe("blog.created", func(msg *nats.Msg) {
			var event events.BlogCreatedEvent
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Println("Error unmarshalling event:", err)
				return
			}

			// Get followers of author
			ctx := context.Background()
			authorID, err := strconv.Atoi(event.AuthorID)
			if err != nil {
				log.Println("Invalid author ID in event:", event.AuthorID)
				return
			}
			followers, err := followerService.GetFollowers(ctx, authorID)
			if err != nil {
				log.Println("Error getting followers:", err)
				return
			}

			//username, _ := userService.GetUsernameByID(authorID) rpc?

			// Notify each follower
			for _, f := range followers {
				notifications.Notify(strconv.Itoa(f.UserID), map[string]interface{}{
					"type":     "new_blog",
					"title":    event.Title,
					"authorId": event.AuthorID,
				})
			}
		})
		if err != nil {
			log.Fatal("Failed to subscribe to blog.created:", err)
		}
		fmt.Println("Subscribed to blog.created events")
	}()

	r := mux.NewRouter()

	r.HandleFunc("/followers/users", followerHandler.CreateUser).Methods("POST")
	r.HandleFunc("/followers/users/{userId}/followers", followerHandler.GetFollowers).Methods("GET")
	r.HandleFunc("/followers/users/{userId}/following", followerHandler.GetFollowing).Methods("GET")
	r.HandleFunc("/followers/users/{userId}/stats", followerHandler.GetUserStats).Methods("GET")
	r.HandleFunc("/followers/users/{userId}/recommendations", followerHandler.GetFollowRecommendations).Methods("GET")

	// Follow/Unfollow routes
	r.HandleFunc("/followers/follow", followerHandler.FollowUser).Methods("POST")
	r.HandleFunc("/followers/unfollow", followerHandler.UnfollowUser).Methods("POST")

	// Check relationships
	r.HandleFunc("/followers/is-following/{followerId}/{followingId}", followerHandler.IsFollowing).Methods("GET")
	r.HandleFunc("/followers/can-comment/{commenterId}/{authorId}", followerHandler.CanUserComment).Methods("GET")

	// WebSocket endpoint
	r.HandleFunc("/followers/ws", wsHandler)

	// Health check
	r.HandleFunc("/followers/health", func(w http.ResponseWriter, r *http.Request) {
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
