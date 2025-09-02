package service

import (
	"blogs/model"
	"blogs/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService struct {
	Repo *repo.BlogRepository
}

func (s *BlogService) GetAllBlogs() ([]model.Blog, error) {
	return s.Repo.FindAll()
}

func (s *BlogService) GetBlog(id primitive.ObjectID) (*model.Blog, error) {
	return s.Repo.FindById(id)
}

func (s *BlogService) CreateBlog(blog *model.Blog) error {
	return s.Repo.Create(blog)
}

func (s *BlogService) GetBlogsByCreator(creatorID int) ([]model.Blog, error) {
	return s.Repo.FindByCreatorID(creatorID)
}
func (s *BlogService) UpdateBlog(blogID primitive.ObjectID, updatedBlog *model.Blog) error {
	return s.Repo.UpdateBlog(blogID, updatedBlog)
}

func (s *BlogService) AddImageToBlog(blogID primitive.ObjectID, imageURL string) error {
	return s.Repo.AddImageToBlog(blogID, imageURL)
}

func (s *BlogService) RemoveImageFromBlog(blogID primitive.ObjectID, imageURL string) error {
	return s.Repo.RemoveImageFromBlog(blogID, imageURL)
}
