package dto

type UserProfile struct {
	Name    string `json:"name" gorm:"not null; type: varchar(20)"`
	Surname string `json:"surname" gorm:"not null; type: varchar(40)"`
	Image   string `json:"image" gorm:"not null"`
	Bio     string `json:"bio" gorm:"not null"`
	Moto    string `json:"moto" gorm:"not null"`
}
