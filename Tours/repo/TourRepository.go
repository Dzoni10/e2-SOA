package repo

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("failed to query tour: %w", err)
	}
	defer cursor.Close(ctx)

	var tours []model.Tour
	if err = cursor.All(ctx, &tours); err != nil {
		fmt.Println("Error decoding tours:", err)
		return nil, fmt.Errorf("failed to decode tours: %w", err)
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

func (r *TourRepository) UpdateStatus(tourID primitive.ObjectID, status model.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"_id": tourID}

	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": now,
		},
	}

	switch status {
	case model.Published:
		update["$set"].(bson.M)["publishedAt"] = now
		update["$unset"] = bson.M{"archivedAt": ""}
	case model.Archived:
		update["$set"].(bson.M)["archivedAt"] = now
	case model.Draft:
		update["$unset"] = bson.M{
			"publishedAt": "",
			"archivedAt":  "",
		}
	}

	_, err := database.TourCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TourRepository) UpdateLengthAndTimes(tourID primitive.ObjectID, length float64, walkingTime, bicycleTime, carTime int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": tourID}
	update := bson.M{
		"$set": bson.M{
			"tourLength":  length,
			"walkingTime": walkingTime,
			"bicycleTime": bicycleTime,
			"carTime":     carTime,
			"updatedAt":   time.Now(),
		},
	}
	_, err := database.TourCollection.UpdateOne(ctx, filter, update)
	return err
}
func (r *TourRepository) UpdateCost(tourID primitive.ObjectID, cost float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": tourID}
	update := bson.M{
		"$set": bson.M{
			"cost":      cost,
			"updatedAt": time.Now(),
		},
	}

	_, err := database.TourCollection.UpdateOne(ctx, filter, update)
	return err
}
