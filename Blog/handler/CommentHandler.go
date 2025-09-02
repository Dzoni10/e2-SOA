package handler

import (
	"blogs/model"
	"blogs/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	Service *service.CommentService
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	blogIDStr := mux.Vars(r)["id"]
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetComments(blogID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
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
