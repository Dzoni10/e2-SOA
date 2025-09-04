package service

import (
	"context"
	"follower-service/model"
	"follower-service/repository"
)

type FollowerService struct {
	Repo *repository.FollowerRepository
}

func NewFollowerService() *FollowerService {
	return &FollowerService{
		Repo: &repository.FollowerRepository{},
	}
}

func (s *FollowerService) CreateUser(ctx context.Context, userID int, name, username string) error {
	return s.Repo.CreateUser(ctx, userID, name, username)
}

func (s *FollowerService) FollowUser(ctx context.Context, followerID, followingID int) error {
	// Check if already following
	isFollowing, err := s.Repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	if isFollowing {
		return nil // Already following, no need to create duplicate relationship
	}

	return s.Repo.FollowUser(ctx, followerID, followingID)
}

func (s *FollowerService) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	return s.Repo.UnfollowUser(ctx, followerID, followingID)
}

func (s *FollowerService) IsFollowing(ctx context.Context, followerID, followingID int) (bool, error) {
	return s.Repo.IsFollowing(ctx, followerID, followingID)
}

func (s *FollowerService) GetFollowers(ctx context.Context, userID int) ([]model.User, error) {
	return s.Repo.GetFollowers(ctx, userID)
}

func (s *FollowerService) GetFollowing(ctx context.Context, userID int) ([]model.User, error) {
	return s.Repo.GetFollowing(ctx, userID)
}

func (s *FollowerService) GetUserStats(ctx context.Context, userID int) (*model.UserFollowStats, error) {
	return s.Repo.GetUserStats(ctx, userID)
}

// CanUserComment checks if a user can comment on another user's blog using userIds
func (s *FollowerService) CanUserComment(ctx context.Context, commenterID, blogAuthorID int) (bool, error) {
	// Users can always comment on their own blogs
	if commenterID == blogAuthorID {
		return true, nil
	}

	// Check if commenter is following the blog author
	return s.Repo.IsFollowing(ctx, commenterID, blogAuthorID)
}

// GetFollowRecommendations returns recommendations for users to follow
func (s *FollowerService) GetFollowRecommendations(ctx context.Context, userID int, limit int) (*model.RecommendationsResponse, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 50 {
		limit = 50 // max limit
	}

	recommendations, err := s.Repo.GetFollowRecommendations(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	return &model.RecommendationsResponse{
		Recommendations: recommendations,
		TotalCount:      len(recommendations),
	}, nil
}
