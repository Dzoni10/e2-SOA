package main

import (
	"blogs/database"
	"blogs/handler"
	"blogs/repo"
	"blogs/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	database.Init()

	// Blogs setup
	blogRepo := &repo.BlogRepository{}
	blogSrv := &service.BlogService{Repo: blogRepo}
	blogHandler := &handler.BlogHandler{Service: blogSrv}

	r := mux.NewRouter()

	// Blogs routes
	r.HandleFunc("/blogs/all", blogHandler.GetAllBlogs).Methods("GET")
	r.HandleFunc("/blogs/{id}", blogHandler.GetBlog).Methods("GET")
	r.HandleFunc("/blogs", blogHandler.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs/creator/{creatorId}", blogHandler.GetBlogsByCreator).Methods("GET")
	r.HandleFunc("/uploads/images/{filename}", blogHandler.ServeImage).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	log.Println("Blog server running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", corsHandler.Handler(r)))
}
