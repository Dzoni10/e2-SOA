package model

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int
type TourDifficulty int

const (
	Draft Status = iota
	Published
	Archived
)

const (
	Easy TourDifficulty = iota
	Medium
	Hard
)

type Tour struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Difficulty  TourDifficulty     `bson:"difficulty" json:"difficulty"`
	Tags        string             `bson:"tags,omitempty" json:"tags"`
	Status      Status             `bson:"status" json:"status"`
	Cost        float64            `bson:"cost" json:"cost"`
	TourLength  float64            `bson:"tourLength" json:"tourLength"`
	CreatorID   int                `bson:"creatorID" json:"creatorID"`
}

func (r Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r TourDifficulty) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r Status) String() string {
	switch r {
	case Draft:
		return "Draft"
	case Published:
		return "Published"
	case Archived:
		return "Archived"
	default:
		return "Unknown"
	}
}

func (r TourDifficulty) String() string {
	switch r {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}
