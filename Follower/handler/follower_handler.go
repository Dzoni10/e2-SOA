package handler

import (
	"context"
	"encoding/json"
	"follower-service/model"
	"follower-service/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FollowerHandler struct {
	Service *service.FollowerService
}

func NewFollowerHandler(service *service.FollowerService) *FollowerHandler {
	return &FollowerHandler{
		Service: service,
	}
}

// CreateUser creates or updates a user
func (h *FollowerHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.Service.CreateUser(ctx, user.UserID, user.Name, user.Username)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.FollowResponse{
		Success: true,
		Message: "User created successfully",
	})
}

// FollowUser creates a follow relationship
func (h *FollowerHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	var req model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FollowerID == req.FollowingID {
		http.Error(w, "User cannot follow themselves", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.Service.FollowUser(ctx, req.FollowerID, req.FollowingID)
	if err != nil {
		log.Printf("Error following user: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.FollowResponse{
		Success:   true,
		Message:   "User followed successfully",
		Following: true,
	})
}

// UnfollowUser removes a follow relationship
func (h *FollowerHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	var req model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := h.Service.UnfollowUser(ctx, req.FollowerID, req.FollowingID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.FollowResponse{
		Success:   true,
		Message:   "User unfollowed successfully",
		Following: false,
	})
}

// IsFollowing checks if one user follows another
func (h *FollowerHandler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	followerID, err := strconv.Atoi(vars["followerId"])
	if err != nil {
		http.Error(w, "Invalid follower ID", http.StatusBadRequest)
		return
	}

	followingID, err := strconv.Atoi(vars["followingId"])
	if err != nil {
		http.Error(w, "Invalid following ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	isFollowing, err := h.Service.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		log.Printf("Error checking follow relationship: %v", err)
		http.Error(w, "Failed to check follow relationship", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"isFollowing": isFollowing,
		"followerId":  followerID,
		"followingId": followingID,
	})
}

// GetFollowers returns list of users who follow the specified user
func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	followers, err := h.Service.GetFollowers(ctx, userID)
	if err != nil {
		log.Printf("Error getting followers: %v", err)
		http.Error(w, "Failed to get followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

// GetFollowing returns list of users that the specified user follows
func (h *FollowerHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	following, err := h.Service.GetFollowing(ctx, userID)
	if err != nil {
		log.Printf("Error getting following: %v", err)
		http.Error(w, "Failed to get following", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

// GetUserStats returns follower/following statistics for a user
func (h *FollowerHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	stats, err := h.Service.GetUserStats(ctx, userID)
	if err != nil {
		log.Printf("Error getting user stats: %v", err)
		http.Error(w, "Failed to get user stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// CanUserComment checks if a user can comment on another user's content
func (h *FollowerHandler) CanUserComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	commenterID, err := strconv.Atoi(vars["commenterId"])
	if err != nil {
		http.Error(w, "Invalid commenter ID", http.StatusBadRequest)
		return
	}

	authorID, err := strconv.Atoi(vars["authorId"])
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	canComment, err := h.Service.CanUserComment(ctx, commenterID, authorID)
	if err != nil {
		log.Printf("Error checking comment permission: %v", err)
		http.Error(w, "Failed to check comment permission", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"canComment":  canComment,
		"commenterId": commenterID,
		"authorId":    authorID,
	})
}

// GetFollowRecommendations returns recommendations for users to follow
func (h *FollowerHandler) GetFollowRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	ctx := context.Background()
	recommendations, err := h.Service.GetFollowRecommendations(ctx, userID, limit)
	if err != nil {
		log.Printf("Error getting follow recommendations: %v", err)
		http.Error(w, "Failed to get follow recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
