package repo

import (
	"blogs/database"
	"blogs/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeRepository struct{}

func (r *LikeRepository) AddLike(like model.Like) error {
	_, err := database.LikeCollection.InsertOne(context.Background(), like)
	return err
}

func (r *LikeRepository) RemoveLike(blogID primitive.ObjectID, userID int) (int64, error) {
	filter := bson.M{"blogId": blogID, "userId": userID}
	result, err := database.LikeCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *LikeRepository) CountLikes(blogID primitive.ObjectID) (int64, error) {
	filter := bson.M{"blogId": blogID}
	return database.LikeCollection.CountDocuments(context.Background(), filter)
}

func (r *LikeRepository) HasUserLikedBlog(blogID primitive.ObjectID, userID int) (bool, error) {
	filter := bson.M{"blogId": blogID, "userId": userID}
	count, err := database.LikeCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
