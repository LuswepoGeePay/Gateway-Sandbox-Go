package config

import (
	"fmt"
	"pg_sandbox/models"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	dsn := "root@tcp(127.0.0.1:3306)/pgsandbox?charset=utf8mb4&parseTime=True&loc=Local"
	//proddsn := "root:password@tcp(127.0.0.1:3306)/hairhunt?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Api{},
		&models.ApiKeys{},
		&models.Notifications{},
		&models.Role{},
		&models.Permission{},
		&models.Transactions{},
		&models.CheckOutUrls{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database connected successfully")
}
