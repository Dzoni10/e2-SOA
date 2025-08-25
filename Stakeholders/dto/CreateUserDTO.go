package dto

import "database-example/model"

type CreateUserDTO struct {
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     model.Role `json:"role"`
}
