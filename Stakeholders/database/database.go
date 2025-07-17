package database

import (
	"database-example/model"
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	//PRE DOKERIZACIJE

	// dsn := "host=localhost user=postgres password=super dbname=Stakeholders sslmode=disable"
	// var err error
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Failed to connect with database", err)
	// }

	// fmt.Println("Database connected")

	// err = DB.AutoMigrate(&model.User{})

	// if err != nil {
	// 	log.Fatal("Failed to migrate database", err)
	// }

	//POSLE DOKERIZACIJE
	// Učitavanje vrednosti iz ENV promenljivih
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Sastavljanje DSN stringa
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect with database:", err)
	}

	fmt.Println("✅ Database connected")

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
