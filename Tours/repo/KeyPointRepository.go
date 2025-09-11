package repo

import (
	"context"
	"tours/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KeyPointRepository struct {
	Collection *mongo.Collection
}

func (r *KeyPointRepository) Create(keyPoint *model.KeyPoint) error {
	_, err := r.Collection.InsertOne(context.TODO(), keyPoint)
	return err
}

func (r *KeyPointRepository) FindByTourID(tourID primitive.ObjectID) ([]model.KeyPoint, error) {
	filter := bson.M{"tourId": tourID}

	// Sort by order field to maintain keypoint sequence
	opts := options.Find().SetSort(bson.D{{"order", 1}})

	cursor, err := r.Collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var keyPoints []model.KeyPoint
	if err = cursor.All(context.TODO(), &keyPoints); err != nil {
		return nil, err
	}

	return keyPoints, nil
}

func (r *KeyPointRepository) FindByID(keyPointID primitive.ObjectID) (*model.KeyPoint, error) {
	filter := bson.M{"_id": keyPointID}

	var keyPoint model.KeyPoint
	err := r.Collection.FindOne(context.TODO(), filter).Decode(&keyPoint)
	if err != nil {
		return nil, err
	}

	return &keyPoint, nil
}

func (r *KeyPointRepository) GetMaxOrderForTour(tourID primitive.ObjectID) (int, error) {
	filter := bson.M{"tourId": tourID}
	opts := options.FindOne().SetSort(bson.D{{"order", -1}})

	var keyPoint model.KeyPoint
	err := r.Collection.FindOne(context.TODO(), filter, opts).Decode(&keyPoint)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil // First keypoint starts at order 1
		}
		return 0, err
	}

	return keyPoint.Order + 1, nil
}

func (r *KeyPointRepository) UpdateTourIdAndOrder(keyPointID, tourID primitive.ObjectID, order int) error {
	filter := bson.M{"_id": keyPointID}
	update := bson.M{
		"$set": bson.M{
			"tourId": tourID,
			"order":  order,
		},
	}

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *KeyPointRepository) BulkUpdateTourId(keyPointIds []primitive.ObjectID, tourID primitive.ObjectID) error {
	// Create bulk write operations
	var writes []mongo.WriteModel

	for i, keyPointId := range keyPointIds {
		filter := bson.M{"_id": keyPointId}
		update := bson.M{
			"$set": bson.M{
				"tourId": tourID,
				"order":  i + 1, // Set order based on array position (1-indexed)
			},
		}

		updateModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update)

		writes = append(writes, updateModel)
	}

	if len(writes) == 0 {
		return nil
	}

	_, err := r.Collection.BulkWrite(context.TODO(), writes)
	return err
}

func (r *KeyPointRepository) UpdateOrder(keyPointID primitive.ObjectID, order int) error {
	filter := bson.M{"_id": keyPointID}
	update := bson.M{
		"$set": bson.M{
			"order": order,
		},
	}

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *KeyPointRepository) BulkReorder(keyPoints []struct {
	ID    primitive.ObjectID `json:"id"`
	Order int                `json:"order"`
}) error {
	var writes []mongo.WriteModel

	for _, kp := range keyPoints {
		filter := bson.M{"_id": kp.ID}
		update := bson.M{
			"$set": bson.M{
				"order": kp.Order,
			},
		}

		updateModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update)

		writes = append(writes, updateModel)
	}

	if len(writes) == 0 {
		return nil
	}

	_, err := r.Collection.BulkWrite(context.TODO(), writes)
	return err
}

func (r *KeyPointRepository) AddImages(keyPointID primitive.ObjectID, images []string) error {
	filter := bson.M{"_id": keyPointID}
	update := bson.M{
		"$push": bson.M{
			"images": bson.M{
				"$each": images,
			},
		},
	}

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *KeyPointRepository) RemoveImage(keyPointID primitive.ObjectID, imagePath string) error {
	filter := bson.M{"_id": keyPointID}
	update := bson.M{
		"$pull": bson.M{
			"images": imagePath,
		},
	}

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *KeyPointRepository) Delete(keyPointID primitive.ObjectID) error {
	filter := bson.M{"_id": keyPointID}
	_, err := r.Collection.DeleteOne(context.TODO(), filter)
	return err
}
