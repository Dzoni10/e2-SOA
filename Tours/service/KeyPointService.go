package service

import (
	"tours/model"
	"tours/repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyPointService struct {
	Repo *repo.KeyPointRepository
}

func (s *KeyPointService) Create(keyPoint *model.KeyPoint) error {
	return s.Repo.Create(keyPoint)
}

func (s *KeyPointService) GetKeyPointsForTour(tourID primitive.ObjectID) ([]model.KeyPoint, error) {
	return s.Repo.FindByTourID(tourID)
}

func (s *KeyPointService) GetByID(keyPointID primitive.ObjectID) (*model.KeyPoint, error) {
	return s.Repo.FindByID(keyPointID)
}

func (s *KeyPointService) GetNextOrderForTour(tourID primitive.ObjectID) (int, error) {
	maxOrder, err := s.Repo.GetMaxOrderForTour(tourID)
	if err != nil {
		return 1, err // Start with 1 if no keypoints exist
	}
	return maxOrder + 1, nil
}

func (s *KeyPointService) UpdateTourIdAndOrder(keyPointID, tourID primitive.ObjectID, order int) error {
	return s.Repo.UpdateTourIdAndOrder(keyPointID, tourID, order)
}

func (s *KeyPointService) BulkUpdateTourId(keyPointIds []primitive.ObjectID, tourID primitive.ObjectID) error {
	return s.Repo.BulkUpdateTourId(keyPointIds, tourID)
}

func (s *KeyPointService) UpdateOrder(keyPointID primitive.ObjectID, order int) error {
	return s.Repo.UpdateOrder(keyPointID, order)
}

func (s *KeyPointService) BulkReorder(keyPoints []struct {
	ID    primitive.ObjectID `json:"id"`
	Order int                `json:"order"`
}) error {
	return s.Repo.BulkReorder(keyPoints)
}

func (s *KeyPointService) AddImages(keyPointID primitive.ObjectID, images []string) error {
	return s.Repo.AddImages(keyPointID, images)
}

func (s *KeyPointService) RemoveImage(keyPointID primitive.ObjectID, imagePath string) error {
	return s.Repo.RemoveImage(keyPointID, imagePath)
}

func (s *KeyPointService) Delete(keyPointID primitive.ObjectID) error {
	return s.Repo.Delete(keyPointID)
}

func (s *KeyPointService) UpdateKeypoint(id primitive.ObjectID, update bson.M) error {
	return s.Repo.UpdateKeypoint(id, update)
}
