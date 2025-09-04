package model

import "time"

// User represents a user node in the graph
type User struct {
	UserID   int    `json:"userId"`
	Username string `json:"username,omitempty"` // Za display u grafu
	Name     string `json:"name,omitempty"`
}

// Follow represents a FOLLOWS relationship between users
type Follow struct {
	FollowerID  int       `json:"followerId"`
	FollowingID int       `json:"followingId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// FollowResponse represents the response for follow operations
type FollowResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Following bool   `json:"following,omitempty"`
}

// UserFollowStats represents follower/following statistics
type UserFollowStats struct {
	UserID         int `json:"userId"`
	FollowersCount int `json:"followersCount"`
	FollowingCount int `json:"followingCount"`
}

// FollowRequest represents a follow/unfollow request
type FollowRequest struct {
	FollowerID  int `json:"followerId"`
	FollowingID int `json:"followingId"`
}

// FollowRecommendation represents a recommendation for following
type FollowRecommendation struct {
	RecommendedUser      User     `json:"recommendedUser"`
	MutualFollowers      []string `json:"mutualFollowers"`
	RecommendationReason string   `json:"recommendationReason"`
}

// RecommendationsResponse represents the response for follow recommendations
type RecommendationsResponse struct {
	Recommendations []FollowRecommendation `json:"recommendations"`
	TotalCount      int                    `json:"totalCount"`
}
