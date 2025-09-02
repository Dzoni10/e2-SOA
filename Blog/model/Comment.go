package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BlogID    primitive.ObjectID `bson:"blogId" json:"blogId"`
	UserID    int                `bson:"userId" json:"userId"`
	Username  string             `bson:"username" json:"username"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
