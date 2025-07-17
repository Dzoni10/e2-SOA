package model

type Role int //Role type for enum

const (
	Admin Role = iota
	Guide
	Tourist
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null; type: varchar(20)"`
	Surname  string `json:"surname" gorm:"not null; type: varchar(40)"`
	Username string `json:"username" gorm:"not null; type: varchar(10)"`
	Password string `json:"password" gorm:"not null; type: varchar(100)"`
	Role     Role   `json:"role" gorm:"not null"`
}
