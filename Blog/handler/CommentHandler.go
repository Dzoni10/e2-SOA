package handler

import (
	"blogs/model"
	"blogs/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	Service     *service.CommentService
	BlogService *service.BlogService
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	blog, err := h.BlogService.GetBlog(blogID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Proverava može li korisnik da vidi komentare
	canComment, err := h.canUserComment(userID, blog.CreatorID)
	if err != nil {
		log.Printf("Error checking comment permission: %v", err)
		http.Error(w, "Unable to verify comment permission", http.StatusInternalServerError)
		return
	}

	if !canComment {
		http.Error(w, "You don't have permission to view these comments. You must follow this user first.", http.StatusForbidden)
		return
	}

	comments, err := h.Service.GetComments(blogID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) canUserComment(commenterID, blogAuthorID int) (bool, error) {
	// Vlasnik bloga uvek može da vidi komentare
	if commenterID == blogAuthorID {
		return true, nil
	}

	followerServiceURL := os.Getenv("FOLLOWER_SERVICE_URL")
	if followerServiceURL == "" {
		followerServiceURL = "http://localhost:8083"
	}

	url := fmt.Sprintf("%s/followers/can-comment/%d/%d", followerServiceURL, commenterID, blogAuthorID)
	log.Printf("Calling follower service: %s", url) // Debug log

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to contact follower service: %v", err)
		return false, fmt.Errorf("failed to contact follower service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Follower service returned status: %d", resp.StatusCode)
		return false, fmt.Errorf("follower service returned status: %d", resp.StatusCode)
	}

	// Kreirajte specifičnu strukturu za odgovor
	var result struct {
		CanComment  bool `json:"canComment"`
		CommenterID int  `json:"commenterId"`
		AuthorID    int  `json:"authorId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode follower service response: %v", err)
		return false, fmt.Errorf("failed to decode follower service response: %v", err)
	}

	log.Printf("Follower service response: canComment=%t, commenterId=%d, authorId=%d",
		result.CanComment, result.CommenterID, result.AuthorID) // Debug log

	return result.CanComment, nil
}

func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// NOVA LOGIKA: Dobijanje blog podataka da znamo ko je autor
	blog, err := h.BlogService.GetBlog(blogID)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// NOVA PROVJERA: Može li korisnik da komentariše?
	canComment, err := h.canUserComment(comment.UserID, blog.CreatorID)
	if err != nil {
		log.Printf("Error checking comment permission: %v", err)
		http.Error(w, "Unable to verify comment permission", http.StatusInternalServerError)
		return
	}

	if !canComment {
		http.Error(w, "You must follow this user to comment on their blog", http.StatusForbidden)
		return
	}

	// Postojeća logika nastavlja...
	comment.ID = primitive.NewObjectID()
	comment.BlogID = blogID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := h.Service.AddComment(&comment); err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}
