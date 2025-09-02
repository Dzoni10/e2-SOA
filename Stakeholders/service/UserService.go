package service

import (
	"database-example/dto"
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

func (service *UserService) GetAllUsersExcept(userID int) ([]model.User, error) {
	return service.Repo.FindAllExcept(userID)
}

func (service *UserService) BlockUser(userID int) (*model.User, error) {
	return service.Repo.BlockUser(userID)
}

func (service *UserService) UpdateUserProfile(u *model.User, dto *dto.UserProfile) {
	if dto.Name != "" {
		u.Name = dto.Name
	}
	if dto.Surname != "" {
		u.Surname = dto.Surname
	}
	if dto.Bio != "" {
		u.Bio = dto.Bio
	}
	if dto.Moto != "" {
		u.Moto = dto.Moto
	}
	if dto.Image != "" {
		u.Image = dto.Image
	}
	service.Repo.DatabaseConnection.Save(&u)
}
