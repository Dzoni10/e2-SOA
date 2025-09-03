package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TourID    primitive.ObjectID `bson:"tourId" json:"tourId"`
	UserID    int                `bson:"userId" json:"userId"`
	Username  string             `bson:"username" json:"username"`
	Rating    int                `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	VisitedAt time.Time          `bson:"visitedAt" json:"visitedAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Images    []string           `bson:"images,omitempty" json:"images"`
}
