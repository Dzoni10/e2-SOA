package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"` // Markdown content
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Images      []string           `bson:"images,omitempty" json:"images"` // Array of image URLs
	CreatorID   int                `bson:"creatorID" json:"creatorID"`
}
