package repo

import (
	"context"
	"time"
	"tours/database"
	"tours/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourExecutionRepository struct{}

func (r *TourExecutionRepository) Create(execution *model.TourExecution) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := database.TourExecutionCollection.InsertOne(ctx, execution)
	return err
}

func (r *TourExecutionRepository) FindById(id primitive.ObjectID) (*model.TourExecution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var execution model.TourExecution
	err := database.TourExecutionCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&execution)
	return &execution, err
}

func (r *TourExecutionRepository) Update(execution *model.TourExecution) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": execution.ID}
	update := bson.M{"$set": execution}
	_, err := database.TourExecutionCollection.UpdateOne(ctx, filter, update)
	return err
}
