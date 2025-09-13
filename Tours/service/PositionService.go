package service

import (
	"tours/model"
	"tours/repo"

	"go.mongodb.org/mongo-driver/bson"
)

type PositionService struct {
	Repo *repo.PositionRepository
}

func (s *PositionService) InitializePosition(position *model.Position) error {
	return s.Repo.InitializePosition(position)
}

func (s *PositionService) UpdatePosition(id int, update bson.M) error {
	return s.Repo.UpdatePosition(id, update)
}

func (s *PositionService) GetPosition(userId int) (*model.Position, error) {
	return s.Repo.FindByUserID(userId)
}
