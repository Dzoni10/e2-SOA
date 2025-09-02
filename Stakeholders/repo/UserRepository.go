package repo

import (
	"database-example/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *UserRepository) FindById(id int) (model.User, error) {
	var user model.User
	dbResult := repo.DatabaseConnection.First(&user, id)

	if dbResult != nil {
		return user, dbResult.Error
	}
	return user, dbResult.Error
}

func (repo *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := repo.DatabaseConnection.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) CreateUser(user *model.User) error {
	dbResult := repo.DatabaseConnection.Create(user)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	print("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *UserRepository) FindAllExcept(userID int) ([]model.User, error) {
	var users []model.User
	result := repo.DatabaseConnection.Where("role != ? AND id != ?", 0, userID).Find(&users)
	return users, result.Error
}

func (repo *UserRepository) BlockUser(userID int) (*model.User, error) {
	var foundUser model.User
	foundUser, err := repo.FindById(userID)
	if err != nil {
		return nil, err
	}
	foundUser.Blocked = true
	dbResult := repo.DatabaseConnection.Save(&foundUser)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &foundUser, nil
}
