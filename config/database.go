package config

import (
	"fmt"
	"os"
	"pg_sandbox/models"
	"pg_sandbox/seeders"

	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// dsn := "root@tcp(127.0.0.1:3306)/pgsandbox?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.ApiKeys{},
		&models.Notifications{},
		&models.Role{},
		&models.Permission{},
		&models.Transactions{},
		&models.CheckOutUrls{},
		&models.APILogs{},
		&models.ActivityLogs{},
		&models.CardUrls{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database connected successfully")

	if err := seeders.SeedRoles(DB); err != nil {
		log.Fatalf("failed to seed roles: %v", err)
	}

	if err := seeders.SeedAdminUser(DB); err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}
}

func LoadEnv() {

	// err := godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

}
