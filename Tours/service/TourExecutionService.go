package service

import (
	"time"
	"tours/model"
	"tours/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourExecutionService struct {
	Repo *repo.TourExecutionRepository
}

func (s *TourExecutionService) GetExecutionById(id primitive.ObjectID) (*model.TourExecution, error) {
	return s.Repo.FindById(id)
}

func (s *TourExecutionService) StartTour(tourID primitive.ObjectID, touristID int, lat, lon float64) (*model.TourExecution, error) {
	execution := &model.TourExecution{
		ID:                 primitive.NewObjectID(),
		TourID:             tourID,
		TouristID:          touristID,
		TourStartDate:      time.Now().UTC(),
		LastActivity:       time.Now().UTC(),
		Status:             model.ExecutionOngoing,
		CompletedPercent:   0,
		CompletedKeyPoints: []model.CompletedKeyPoint{},
		CurrentLatitude:    lat,
		CurrentLongitude:   lon,
	}
	err := s.Repo.Create(execution)
	return execution, err
}

func (s *TourExecutionService) FinishTour(id primitive.ObjectID) error {
	execution, err := s.Repo.FindById(id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	execution.Status = model.ExecutionCompleted
	execution.TourEndDate = &now
	execution.LastActivity = now
	return s.Repo.Update(execution)
}

func (s *TourExecutionService) AbandonTour(id primitive.ObjectID) error {
	execution, err := s.Repo.FindById(id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	execution.Status = model.ExecutionAbandoned
	execution.TourEndDate = &now
	execution.LastActivity = now
	return s.Repo.Update(execution)
}

func (s *TourExecutionService) AddCompletedKeyPoint(executionID primitive.ObjectID, keyPointID primitive.ObjectID, totalKeyPoints int) error {
	execution, err := s.Repo.FindById(executionID)
	if err != nil {
		return err
	}
	execution.CompletedKeyPoints = append(execution.CompletedKeyPoints, model.CompletedKeyPoint{
		KeyPointID: keyPointID,
		ReachedAt:  time.Now().UTC(),
	})
	execution.LastActivity = time.Now().UTC()
	execution.CompletedPercent = float64(len(execution.CompletedKeyPoints)) / float64(totalKeyPoints) * 100
	return s.Repo.Update(execution)
}

func (s *TourExecutionService) UpdateLocation(executionID primitive.ObjectID, lat, lon float64) error {
	execution, err := s.Repo.FindById(executionID)
	if err != nil {
		return err
	}

	// update polja
	execution.CurrentLatitude = lat
	execution.CurrentLongitude = lon
	execution.LastActivity = time.Now().UTC()

	// sacuvaj promene
	return s.Repo.Update(execution)
}
