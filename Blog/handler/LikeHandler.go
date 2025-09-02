package handler

import (
	"blogs/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeHandler struct {
	Service *service.LikeService
}

// POST /blogs/{id}/likes
func (h *LikeHandler) LikeBlog(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	var req struct {
		UserID int `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	err = h.Service.LikeBlog(blogObjID, req.UserID)
	if errors.Is(err, service.ErrAlreadyLiked) {
		http.Error(w, "Already liked", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Failed to like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Liked successfully"})
}

// DELETE /blogs/{id}/likes
func (h *LikeHandler) UnlikeBlog(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	var req struct {
		UserID int `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	deleted, err := h.Service.UnlikeBlog(blogObjID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to unlike", http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Unliked successfully"})
}

// GET /blogs/{id}/likes/count
func (h *LikeHandler) CountLikes(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	count, err := h.Service.CountLikes(blogObjID)
	if err != nil {
		http.Error(w, "Failed to count likes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"count": count})
}

// GET /blogs/{id}/likes/{userId}
func (h *LikeHandler) HasUserLiked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID := vars["id"]
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	liked, err := h.Service.HasUserLiked(blogObjID, userID)
	if err != nil {
		http.Error(w, "Failed to check like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}
