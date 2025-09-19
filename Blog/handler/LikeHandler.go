package handler

import (
	"blogs/logger"
	"blogs/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
		logger.Error("Failed to decode like request", err, logrus.Fields{
			"blog_id": blogID,
			"action":  "like_blog",
		})
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	logger.Info("User attempting to like blog", logrus.Fields{
		"user_id": req.UserID,
		"blog_id": blogID,
		"action":  "like_blog",
	})

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		logger.Error("Invalid blog ID format", err, logrus.Fields{
			"blog_id": blogID,
			"user_id": req.UserID,
			"action":  "like_blog",
		})
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	err = h.Service.LikeBlog(blogObjID, req.UserID)
	if errors.Is(err, service.ErrAlreadyLiked) {
		logger.Warn("User tried to like already liked blog", logrus.Fields{
			"user_id": req.UserID,
			"blog_id": blogID,
			"action":  "like_blog",
		})
		http.Error(w, "Already liked", http.StatusConflict)
		return
	}
	if err != nil {
		logger.Error("Failed to like blog", err, logrus.Fields{
			"user_id": req.UserID,
			"blog_id": blogID,
			"action":  "like_blog",
		})
		http.Error(w, "Failed to like", http.StatusInternalServerError)
		return
	}

	logger.Info("User successfully liked blog", logrus.Fields{
		"user_id": req.UserID,
		"blog_id": blogID,
		"action":  "like_blog",
	})

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
		logger.Error("Failed to decode unlike request", err, logrus.Fields{
			"blog_id": blogID,
			"action":  "unlike_blog",
		})
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	logger.Info("User attempting to unlike blog", logrus.Fields{
		"user_id": req.UserID,
		"blog_id": blogID,
		"action":  "unlike_blog",
	})

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		logger.Error("Invalid blog ID format for unlike", err, logrus.Fields{
			"blog_id": blogID,
			"user_id": req.UserID,
			"action":  "unlike_blog",
		})
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	deleted, err := h.Service.UnlikeBlog(blogObjID, req.UserID)
	if err != nil {
		logger.Error("Failed to unlike blog", err, logrus.Fields{
			"user_id": req.UserID,
			"blog_id": blogID,
			"action":  "unlike_blog",
		})
		http.Error(w, "Failed to unlike", http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		logger.Warn("User tried to unlike non-existing like", logrus.Fields{
			"user_id": req.UserID,
			"blog_id": blogID,
			"action":  "unlike_blog",
		})
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}

	logger.Info("User successfully unliked blog", logrus.Fields{
		"user_id": req.UserID,
		"blog_id": blogID,
		"action":  "unlike_blog",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Unliked successfully"})
}

// GET /blogs/{id}/likes/count
func (h *LikeHandler) CountLikes(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	logger.Debug("Counting likes for blog", logrus.Fields{
		"blog_id": blogID,
		"action":  "count_likes",
	})

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		logger.Error("Invalid blog ID format for count", err, logrus.Fields{
			"blog_id": blogID,
			"action":  "count_likes",
		})
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	count, err := h.Service.CountLikes(blogObjID)
	if err != nil {
		logger.Error("Failed to count likes", err, logrus.Fields{
			"blog_id": blogID,
			"action":  "count_likes",
		})
		http.Error(w, "Failed to count likes", http.StatusInternalServerError)
		return
	}

	logger.Debug("Successfully counted likes", logrus.Fields{
		"blog_id": blogID,
		"count":   count,
		"action":  "count_likes",
	})

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
		logger.Error("Invalid user ID format", err, logrus.Fields{
			"blog_id":     blogID,
			"user_id_str": userIDStr,
			"action":      "check_user_liked",
		})
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	logger.Debug("Checking if user liked blog", logrus.Fields{
		"user_id": userID,
		"blog_id": blogID,
		"action":  "check_user_liked",
	})

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		logger.Error("Invalid blog ID format for user like check", err, logrus.Fields{
			"blog_id": blogID,
			"user_id": userID,
			"action":  "check_user_liked",
		})
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	liked, err := h.Service.HasUserLiked(blogObjID, userID)
	if err != nil {
		logger.Error("Failed to check if user liked blog", err, logrus.Fields{
			"user_id": userID,
			"blog_id": blogID,
			"action":  "check_user_liked",
		})
		http.Error(w, "Failed to check like", http.StatusInternalServerError)
		return
	}

	logger.Debug("User like status checked", logrus.Fields{
		"user_id": userID,
		"blog_id": blogID,
		"liked":   liked,
		"action":  "check_user_liked",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}
