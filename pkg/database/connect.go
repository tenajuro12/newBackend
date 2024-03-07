package database

import (
	models2 "github.com/tenajuro12/newBackend/pkg/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	dsn := os.Getenv("DSN")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	} else {
		log.Println("connect successfully")
	}
	DB = database
	database.AutoMigrate(
		&models2.User{},
		&models2.Blog{},
		&models2.Comment{},
		&models2.Follow{},
	)

}
