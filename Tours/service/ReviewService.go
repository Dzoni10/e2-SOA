package service

import (
	"time"
	"tours/model"
	"tours/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	Repo *repo.ReviewRepository
}

func (s *ReviewService) Create(review *model.Review) error {
	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	if review.VisitedAt.IsZero() {
		review.VisitedAt = review.CreatedAt
	}
	return s.Repo.Create(review)
}

func (s *ReviewService) GetReviewsForTour(tourID primitive.ObjectID) ([]model.Review, error) {
	return s.Repo.FindByTourID(tourID)
}

func (s *ReviewService) AddImages(reviewID primitive.ObjectID, images []string) error {
	return s.Repo.AddImages(reviewID, images)
}

func (s *ReviewService) RemoveImage(reviewID primitive.ObjectID, imagePath string) error {
	return s.Repo.RemoveImage(reviewID, imagePath)
}
