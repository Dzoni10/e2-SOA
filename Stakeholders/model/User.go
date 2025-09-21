package model

import (
	"encoding/json"
)

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
	Password string `json:"-" gorm:"not null; type: varchar(100)"` // OVDE MINUS ZNACI DA SE NE SALJE KA FRONTU UD ODGOVORUI
	Role     Role   `json:"role" gorm:"not null"`
	Blocked  bool   `json:"blocked" gorm:"not null"`
	Image    string `json:"image" gorm:"not null"`
	Bio      string `json:"bio" gorm:"not null"`
	Moto     string `json:"moto" gorm:"not null"`
}

// Konverzija enuma u string
func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r Role) String() string {
	switch r {
	case Admin:
		return "Admin"
	case Guide:
		return "Guide"
	case Tourist:
		return "Tourist"
	default:
		return "Unknown"
	}
}
