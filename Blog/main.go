package main

import (
	"blogs/database"
	"blogs/handler"
	"blogs/logger"
	"blogs/repo"
	"blogs/service"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {

	logger.Init("blog-service")

	logger.Info("Starting blog service", logrus.Fields{
		"port":    "8082",
		"version": "1.0.0",
	})

	database.Init()
	logger.Info("Database initialized successfully")

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
	r.Use(loggingMiddleware)

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
	logger.Info("Blog server running", logrus.Fields{"port": "8082"})

	if err := http.ListenAndServe(":8082", r); err != nil {
		logger.Error("Server failed to start", err)
		os.Exit(1)

	}
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("HTTP Request", logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})
		next.ServeHTTP(w, r)
	})
}
