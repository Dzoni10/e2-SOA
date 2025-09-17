package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExecutionStatus string

const (
	ExecutionCompleted ExecutionStatus = "COMPLETED"
	ExecutionAbandoned ExecutionStatus = "ABANDONED"
	ExecutionOngoing   ExecutionStatus = "ONGOING"
)

type CompletedKeyPoint struct {
	KeyPointID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReachedAt  time.Time          `bson:"reachedAt" json:"reachedAt"`
}

type TourExecution struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TourID             primitive.ObjectID  `bson:"tourId" json:"tourId"`
	TouristID          int                 `bson:"touristId" json:"touristId"`
	TourStartDate      time.Time           `bson:"tourStartDate" json:"tourStartDate"`
	TourEndDate        *time.Time          `bson:"tourEndDate,omitempty" json:"tourEndDate,omitempty"`
	LastActivity       time.Time           `bson:"lastActivity" json:"lastActivity"`
	Status             ExecutionStatus     `bson:"status" json:"status"`
	CompletedPercent   float64             `bson:"completedPercent" json:"completedPercent"`
	CompletedKeyPoints []CompletedKeyPoint `bson:"completedKeyPoints" json:"completedKeyPoints"`
	CurrentLatitude    float64             `bson:"currentLatitude" json:"currentLatitude"`
	CurrentLongitude   float64             `bson:"currentLongitude" json:"currentLongitude"`
}
