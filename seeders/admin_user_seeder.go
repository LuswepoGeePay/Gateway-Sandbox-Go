package seeders

import (
	"errors"
	"fmt"
	"pg_sandbox/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) error {
	var adminRole models.Role

	// Ensure admin role exists
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return errors.New("admin role not found, run role seeder first")
	}

	// Check if admin already exists
	var existing models.User
	err := db.Where("email = ?", "luswepo@geepay.co.zm").First(&existing).Error
	if err == nil {
		fmt.Println("Admin user already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password@123"), bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	admin := models.User{
		ID:            uuid.New(),
		Fullname:      "Luswepo Silumbwe",
		Email:         "luswepo@geepay.co.zm",
		Password:      string(hashedPassword),
		Phone:         "0000000000",
		RoleID:        adminRole.ID.String(),
		Status:        "active",
		EmailVerified: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	fmt.Println("Seeded admin user: luswepo@geepay.co.zm")
	return nil
}
