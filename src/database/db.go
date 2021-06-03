package database

import (
	"ambassador/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dsn := "host=db user=root password=root dbname=ambassador port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect with")
	}
}
func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Product{}, models.Link{}, models.Order{}, models.OrderItem{})
}
