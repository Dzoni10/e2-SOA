package service

import (
	"tours/model"
	"tours/repo"
	"tours/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourService struct {
	Repo         *repo.TourRepository
	KeyPointRepo *repo.KeyPointRepository
}

func (s *TourService) GetAllTours() ([]model.Tour, error) {
	return s.Repo.FindAll()
}

func (s *TourService) GetTour(id primitive.ObjectID) (*model.Tour, error) {
	return s.Repo.FindById(id)
}

func (s *TourService) CreateTour(tour *model.Tour) error {
	return s.Repo.Create(tour)
}

func (s *TourService) UpdateTourLength(tourID primitive.ObjectID) error {
	keypoints, err := s.KeyPointRepo.FindByTourID(tourID)
	if err != nil {
		return err
	}

	length := utils.CalculateTourLength(keypoints)

	return s.Repo.UpdateLength(tourID, length)
}
