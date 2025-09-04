package repository

import (
	"context"
	"fmt"

	"follower-service/database"
	"follower-service/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct{}

// CreateUser creates or updates a user node in Neo4j using userId as primary key
func (r *FollowerRepository) CreateUser(ctx context.Context, userID int, name, username string) error {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		MERGE (u:User {userId: $userId})
		SET u.name = $name, u.username = $username, u.updatedAt = datetime()
		RETURN u
	`

	params := map[string]interface{}{
		"userId":   userID,
		"name":     name,
		"username": username,
	}

	_, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create/update user: %v", err)
	}

	return nil
}

// FollowUser creates a FOLLOWS relationship between two users using userIds
func (r *FollowerRepository) FollowUser(ctx context.Context, followerID, followingID int) error {
	if followerID == followingID {
		return fmt.Errorf("user cannot follow themselves")
	}

	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		MERGE (follower:User {userId: $followerId})
		MERGE (following:User {userId: $followingId})
		MERGE (follower)-[r:FOLLOWS]->(following)
		SET r.createdAt = datetime()
		RETURN r
	`

	params := map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	}

	_, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create follow relationship: %v", err)
	}

	return nil
}

// UnfollowUser removes the FOLLOWS relationship between two users using userIds
func (r *FollowerRepository) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {userId: $followerId})-[r:FOLLOWS]->(following:User {userId: $followingId})
		DELETE r
		RETURN count(r) as deletedCount
	`

	params := map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to remove follow relationship: %v", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return fmt.Errorf("relationship not found")
	}

	deletedCount, _ := record.Get("deletedCount")
	if deletedCount.(int64) == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	return nil
}

// IsFollowing checks if followerID is following followingID
func (r *FollowerRepository) IsFollowing(ctx context.Context, followerID, followingID int) (bool, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {userId: $followerId})-[r:FOLLOWS]->(following:User {userId: $followingId})
		RETURN count(r) > 0 as isFollowing
	`

	params := map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return false, fmt.Errorf("failed to check follow relationship: %v", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		isFollowing, _ := record.Get("isFollowing")
		return isFollowing.(bool), nil
	}

	return false, nil
}

// GetFollowers returns list of users who follow the given user
func (r *FollowerRepository) GetFollowers(ctx context.Context, userID int) ([]model.User, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User)-[r:FOLLOWS]->(user:User {userId: $userId})
		RETURN follower.userId as userId, follower.username as username, follower.name as name
		ORDER BY r.createdAt DESC
	`

	params := map[string]interface{}{
		"userId": userID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %v", err)
	}

	var followers []model.User
	for result.Next(ctx) {
		record := result.Record()
		userIdValue, _ := record.Get("userId")
		usernameValue, _ := record.Get("username")
		nameValue, _ := record.Get("name")

		var name, username string
		if nameValue != nil {
			name = nameValue.(string)
		}
		if usernameValue != nil {
			username = usernameValue.(string)
		}

		followers = append(followers, model.User{
			UserID:   int(userIdValue.(int64)),
			Username: username,
			Name:     name,
		})
	}

	return followers, nil
}

// GetFollowing returns list of users that the given user follows
func (r *FollowerRepository) GetFollowing(ctx context.Context, userID int) ([]model.User, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (user:User {userId: $userId})-[r:FOLLOWS]->(following:User)
		RETURN following.userId as userId, following.username as username, following.name as name
		ORDER BY r.createdAt DESC
	`

	params := map[string]interface{}{
		"userId": userID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %v", err)
	}

	var following []model.User
	for result.Next(ctx) {
		record := result.Record()
		userIdValue, _ := record.Get("userId")
		usernameValue, _ := record.Get("username")
		nameValue, _ := record.Get("name")

		var name, username string
		if nameValue != nil {
			name = nameValue.(string)
		}
		if usernameValue != nil {
			username = usernameValue.(string)
		}

		following = append(following, model.User{
			UserID:   int(userIdValue.(int64)),
			Username: username,
			Name:     name,
		})
	}

	return following, nil
}

// GetUserStats returns follower/following counts for a user
func (r *FollowerRepository) GetUserStats(ctx context.Context, userID int) (*model.UserFollowStats, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (user:User {userId: $userId})
		OPTIONAL MATCH (follower:User)-[:FOLLOWS]->(user)
		OPTIONAL MATCH (user)-[:FOLLOWS]->(following:User)
		RETURN 
			user.userId as userId,
			count(DISTINCT follower) as followersCount,
			count(DISTINCT following) as followingCount
	`

	params := map[string]interface{}{
		"userId": userID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %v", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		userIdValue, _ := record.Get("userId")
		followersCountValue, _ := record.Get("followersCount")
		followingCountValue, _ := record.Get("followingCount")

		return &model.UserFollowStats{
			UserID:         int(userIdValue.(int64)),
			FollowersCount: int(followersCountValue.(int64)),
			FollowingCount: int(followingCountValue.(int64)),
		}, nil
	}

	return &model.UserFollowStats{
		UserID:         userID,
		FollowersCount: 0,
		FollowingCount: 0,
	}, nil
}

// GetFollowRecommendations returns recommendations based on userId
func (r *FollowerRepository) GetFollowRecommendations(ctx context.Context, userID int, limit int) ([]model.FollowRecommendation, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		// Pronađi korisnike koje prate oni koje ja pratim, a ja ih ne pratim
		MATCH (me:User {userId: $userId})-[:FOLLOWS]->(following:User)-[:FOLLOWS]->(recommended:User)
		WHERE NOT (me)-[:FOLLOWS]->(recommended) AND recommended.userId <> $userId
		
		// Grupiši po preporučenom korisniku i prebrojaj zajedničke pratioce
		WITH recommended, collect(DISTINCT following.username) as mutualFollowers
		WHERE size(mutualFollowers) > 0
		
		RETURN 
			recommended.userId as userId,
			recommended.username as username, 
			recommended.name as name,
			mutualFollowers,
			size(mutualFollowers) as mutualCount
		ORDER BY mutualCount DESC, recommended.username
		LIMIT $limit
	`

	params := map[string]interface{}{
		"userId": userID,
		"limit":  limit,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get follow recommendations: %v", err)
	}

	var recommendations []model.FollowRecommendation
	for result.Next(ctx) {
		record := result.Record()

		userIdValue, _ := record.Get("userId")
		usernameValue, _ := record.Get("username")
		nameValue, _ := record.Get("name")
		mutualFollowersValue, _ := record.Get("mutualFollowers")
		mutualCountValue, _ := record.Get("mutualCount")

		var recUserID int
		var recUsername, name string

		if userIdValue != nil {
			recUserID = int(userIdValue.(int64))
		}
		if usernameValue != nil {
			recUsername = usernameValue.(string)
		}
		if nameValue != nil {
			name = nameValue.(string)
		}

		mutualFollowers := []string{}
		if mutualFollowersValue != nil {
			mutualList := mutualFollowersValue.([]interface{})
			for _, mutual := range mutualList {
				if mutualStr, ok := mutual.(string); ok {
					mutualFollowers = append(mutualFollowers, mutualStr)
				}
			}
		}

		mutualCount := int(mutualCountValue.(int64))
		reason := fmt.Sprintf("Prate ih %d korisnika koje i ti pratiš", mutualCount)

		recommendations = append(recommendations, model.FollowRecommendation{
			RecommendedUser: model.User{
				UserID:   recUserID,
				Username: recUsername,
				Name:     name,
			},
			MutualFollowers:      mutualFollowers,
			RecommendationReason: reason,
		})
	}

	return recommendations, nil
}
