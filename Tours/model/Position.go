package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Position struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    int                `bson:"userId" json:"userId"`
	Latitude  float64            `bson:"latitude" json:"latitude"`
	Longitude float64            `bson:"longitude" json:"longitude"`
}
