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

func (s *TourService) UpdateTourStatus(tourID primitive.ObjectID, status model.Status) error {
	return s.Repo.UpdateStatus(tourID, status)
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

// Nova helper funkcija
func (s *TourService) CalculateTourMetrics(keypoints []model.KeyPoint) (float64, int, int, int) {
	if len(keypoints) < 2 {
		return 0, 0, 0, 0
	}

	length := utils.CalculateTourLength(keypoints)

	walkingTime := int((length / 5.0) * 60)  // 5 km/h peške
	bicycleTime := int((length / 15.0) * 60) // 15 km/h bicikl
	carTime := int((length / 50.0) * 60)     // 50 km/h auto u gradu

	return length, walkingTime, bicycleTime, carTime
}

// Izmeni postojeću funkciju da koristi helper
func (s *TourService) UpdateTourLength(tourID primitive.ObjectID) error {
	keypoints, err := s.KeyPointRepo.FindByTourID(tourID)
	if err != nil {
		return err
	}

	length, walkingTime, bicycleTime, carTime := s.CalculateTourMetrics(keypoints)

	return s.Repo.UpdateLengthAndTimes(tourID, length, walkingTime, bicycleTime, carTime)
}

func (s *TourService) UpdateTourCost(tourID primitive.ObjectID, cost float64) error {
	return s.Repo.UpdateCost(tourID, cost)
}
