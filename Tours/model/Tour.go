package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int
type TourDifficulty int
type TransportType int

const (
	Walking TransportType = iota
	Bicycle
	Car
)

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
	Tags        []string           `bson:"tags,omitempty" json:"tags"`
	Status      Status             `bson:"status" json:"status"`
	Cost        float64            `bson:"cost" json:"cost"`
	TourLength  float64            `bson:"tourLength" json:"tourLength"`
	CreatorID   int                `bson:"creatorID" json:"creatorID"`
	PublishedAt *time.Time         `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	ArchivedAt  *time.Time         `bson:"archivedAt,omitempty" json:"archivedAt,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	WalkingTime int                `bson:"walkingTime" json:"walkingTime"`
	BicycleTime int                `bson:"bicycleTime" json:"bicycleTime"`
	CarTime     int                `bson:"carTime" json:"carTime"` //SVE JE U MINUTIMA PA CU NA FRONTU DELITI SA 60 AKO TREBA U SATIMA

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
func (t TransportType) String() string {
	switch t {
	case Walking:
		return "Walking"
	case Bicycle:
		return "Bicycle"
	case Car:
		return "Car"
	default:
		return "Unknown"
	}
}

func (t TransportType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
