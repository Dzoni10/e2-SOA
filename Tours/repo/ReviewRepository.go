package repo

import (
	"context"
	"time"
	"tours/database"
	"tours/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewRepository struct{}

func (r *ReviewRepository) Create(review *model.Review) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.ReviewCollection.InsertOne(ctx, review)
	return err
}

func (r *ReviewRepository) FindByTourID(tourID primitive.ObjectID) ([]model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.ReviewCollection.Find(ctx, bson.M{"tourId": tourID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []model.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}
func (r *ReviewRepository) AddImages(reviewID primitive.ObjectID, images []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"images": bson.M{"$each": images}}}
	_, err := database.ReviewCollection.UpdateOne(ctx, bson.M{"_id": reviewID}, update)
	return err
}

func (r *ReviewRepository) RemoveImage(reviewID primitive.ObjectID, imagePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$pull": bson.M{"images": imagePath}}
	_, err := database.ReviewCollection.UpdateOne(ctx, bson.M{"_id": reviewID}, update)
	return err
}
