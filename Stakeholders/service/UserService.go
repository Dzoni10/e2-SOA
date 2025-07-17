package service

import (
	"database-example/model"
	"database-example/repo"
	"errors"
)

type UserService struct {
	Repo *repo.UserRepository
}

func (u *UserService) FindUser(id int) (model.User, error) {
	return u.Repo.FindById(id)
}

func (u *UserService) Create(user *model.User) error {
	return u.Repo.CreateUser(user)
}

func (s *UserService) Authenticate(username, password string) (*model.User, error) {
	user, err := s.Repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Password != password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (service *UserService) GetAllUsers() ([]model.User, error) {
	return service.Repo.FindAll()
}
