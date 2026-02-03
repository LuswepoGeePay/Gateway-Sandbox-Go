package seeders

import (
	"fmt"
	"pg_sandbox/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{
			ID:   uuid.New(),
			Name: "admin",
		},
		{
			ID:   uuid.New(),
			Name: "merchant",
		},
		{
			ID:   uuid.New(),
			Name: "independent",
		},
	}

	for _, role := range roles {
		var existing models.Role

		err := db.Where("name = ?", role.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
			fmt.Println("Seeded role:", role.Name)
		}
	}

	return nil
}
