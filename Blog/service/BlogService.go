package service

import (
	"blogs/model"
	"blogs/repo"
	"local/common/saga/events"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService struct {
	Repo      *repo.BlogRepository
	Publisher *events.Publisher
}

func (s *BlogService) GetAllBlogs() ([]model.Blog, error) {
	return s.Repo.FindAll()
}

func (s *BlogService) GetBlog(id primitive.ObjectID) (*model.Blog, error) {
	return s.Repo.FindById(id)
}

func (s *BlogService) CreateBlog(blog *model.Blog) error {
	if err := s.Repo.Create(blog); err != nil {
		return err
	}

	// Objavi event na NATS
	if s.Publisher != nil {
		event := events.BlogCreatedEvent{
			BlogID:    blog.ID.Hex(),
			AuthorID:  strconv.Itoa(blog.CreatorID),
			Title:     blog.Title,
			Content:   blog.Description,
			CreatedAt: blog.CreatedAt.Format(time.RFC3339),
		}
		_ = s.Publisher.Publish("blog.created", event)
	}

	return nil
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
