package main

import (
	"blogs/database"
	"blogs/handler"
	"blogs/repo"
	"blogs/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	database.Init()

	// Blogs setup
	blogRepo := &repo.BlogRepository{}
	blogSrv := &service.BlogService{Repo: blogRepo}
	blogHandler := &handler.BlogHandler{Service: blogSrv}
	commentRepo := &repo.CommentRepository{}
	commentSrv := &service.CommentService{Repo: commentRepo}
	commentHandler := &handler.CommentHandler{
		Service:     commentSrv,
		BlogService: blogSrv,
	}

	likeSrv := &service.LikeService{}
	likeHandler := &handler.LikeHandler{Service: likeSrv}

	r := mux.NewRouter()

	// Blogs routes
	r.HandleFunc("/blogs/all", blogHandler.GetAllBlogs).Methods("GET")
	r.HandleFunc("/blogs/{id}", blogHandler.GetBlog).Methods("GET")
	r.HandleFunc("/blogs", blogHandler.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs/creator/{creatorId}", blogHandler.GetBlogsByCreator).Methods("GET")
	r.HandleFunc("/uploads/images/{filename}", blogHandler.ServeImage).Methods("GET")
	r.HandleFunc("/blogs/{id}", blogHandler.UpdateBlog).Methods("PUT")
	r.HandleFunc("/blogs/{id}/add-image", blogHandler.AddImageToBlog).Methods("POST")
	r.HandleFunc("/blogs/{id}/remove-image", blogHandler.RemoveImageFromBlog).Methods("DELETE")

	r.HandleFunc("/blogs/{id}/comments", commentHandler.GetComments).Methods("GET")
	r.HandleFunc("/blogs/{id}/comments", commentHandler.AddComment).Methods("POST")
	r.HandleFunc("/blogs/{id}/likes", likeHandler.LikeBlog).Methods("POST")
	r.HandleFunc("/blogs/{id}/likes", likeHandler.UnlikeBlog).Methods("DELETE")
	r.HandleFunc("/blogs/{id}/likes/count", likeHandler.CountLikes).Methods("GET")
	r.HandleFunc("/blogs/{id}/likes/{userId}", likeHandler.HasUserLiked).Methods("GET")

	/*corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	*/
	log.Println("Blog server running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
