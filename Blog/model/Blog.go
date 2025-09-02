package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	UserID    int       `bson:"userId" json:"userId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"` // Markdown content
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Images      []string           `bson:"images,omitempty" json:"images"` // Array of image URLs
	CreatorID   int                `bson:"creatorID" json:"creatorID"`
	Likes       []Like             `bson:"likes,omitempty" json:"likes,omitempty"`
}
