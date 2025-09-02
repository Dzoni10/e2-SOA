package repo

import (
	"blogs/database"
	"blogs/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentRepository struct{}

func (r *CommentRepository) FindByBlogID(blogID primitive.ObjectID) ([]model.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.CommentCollection.Find(ctx, bson.M{"blogId": blogID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err := database.CommentCollection.InsertOne(ctx, comment)
	return err
}
