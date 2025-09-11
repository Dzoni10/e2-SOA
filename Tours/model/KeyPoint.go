package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type KeyPoint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TourID      primitive.ObjectID `bson:"tourId" json:"tourId"`
	Latitude    float64            `bson:"latitude" json:"latitude"`
	Longitude   float64            `bson:"longitude" json:"longitude"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"desc" json:"desc"`
	Order       int                `bson:"order" json:"order"`
	Images      []string           `bson:"images,omitempty" json:"images"`
}
