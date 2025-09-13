package repo

import (
	"context"
	"time"
	"tours/database"
	"tours/model"

	"go.mongodb.org/mongo-driver/bson"
)

type PositionRepository struct{}

func (r *PositionRepository) InitializePosition(position *model.Position) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.PositionCollection.InsertOne(ctx, position)
	return err
}

func (r *PositionRepository) UpdatePosition(userId int, update bson.M) error {
	_, err := database.PositionCollection.UpdateOne(context.TODO(),
		bson.M{"userId": userId},
		update,
	)
	return err
}

func (r *PositionRepository) FindByUserID(userId int) (*model.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var position model.Position
	err := database.PositionCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(&position)

	if err != nil {
		return nil, err
	}
	return &position, nil
}
