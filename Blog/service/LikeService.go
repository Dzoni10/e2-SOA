package service

import (
	"blogs/database"
	"blogs/model"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrAlreadyLiked = errors.New("already liked")

type LikeService struct{}

func (s *LikeService) LikeBlog(blogID primitive.ObjectID, userID int) error {
	ctx := context.Background()

	// provera da li user već lajkovao
	filter := bson.M{"_id": blogID, "likes.userId": userID}
	count, err := database.BlogCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrAlreadyLiked
	}

	// dodaj novi like
	update := bson.M{
		"$push": bson.M{
			"likes": model.Like{
				UserID:    userID,
				CreatedAt: time.Now(),
			},
		},
	}

	_, err = database.BlogCollection.UpdateByID(ctx, blogID, update)
	return err
}

func (s *LikeService) UnlikeBlog(blogID primitive.ObjectID, userID int) (int64, error) {
	ctx := context.Background()

	update := bson.M{
		"$pull": bson.M{
			"likes": bson.M{"userId": userID},
		},
	}

	result, err := database.BlogCollection.UpdateByID(ctx, blogID, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (s *LikeService) CountLikes(blogID primitive.ObjectID) (int64, error) {
	ctx := context.Background()
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: blogID}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "count", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$likes", bson.A{}}},
				}},
			}},
		}}},
	}

	cursor, err := database.BlogCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) > 0 {
		switch v := result[0]["count"].(type) {
		case int32:
			return int64(v), nil
		case int64:
			return v, nil
		default:
			return 0, fmt.Errorf("unexpected type for count: %T", v)
		}
	}
	return 0, nil
}

func (s *LikeService) HasUserLiked(blogID primitive.ObjectID, userID int) (bool, error) {
	ctx := context.Background()
	filter := bson.M{"_id": blogID, "likes.userId": userID}
	count, err := database.BlogCollection.CountDocuments(ctx, filter)
	return count > 0, err
}
