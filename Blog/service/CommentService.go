package service

import (
	"blogs/model"
	"blogs/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	Repo *repo.CommentRepository
}

func (s *CommentService) GetComments(blogID primitive.ObjectID) ([]model.Comment, error) {
	return s.Repo.FindByBlogID(blogID)
}

func (s *CommentService) AddComment(comment *model.Comment) error {
	return s.Repo.Create(comment)
}
