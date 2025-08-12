package service

import (
	"tours/model"
	"tours/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourService struct {
	Repo *repo.TourRepository
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
