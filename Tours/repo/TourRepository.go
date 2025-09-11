package repo

import (
	"context"
	"time"
	"tours/database"
	"tours/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourRepository struct{}

func (r *TourRepository) FindAll() ([]model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.TourCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tours []model.Tour
	if err = cursor.All(ctx, &tours); err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *TourRepository) FindById(id primitive.ObjectID) (*model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tour model.Tour
	err := database.TourCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tour)

	if err != nil {
		return nil, err
	}
	return &tour, nil
}

func (r *TourRepository) Create(tour *model.Tour) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.TourCollection.InsertOne(ctx, tour)
	return err
}

func (r *TourRepository) UpdateLength(tourID primitive.ObjectID, length float64) error {
	filter := bson.M{"_id": tourID}
	update := bson.M{"$set": bson.M{"tourLength": length}}
	_, err := database.TourCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
