package credentialsservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/api"
	"pg_sandbox/utils"

	"gorm.io/gorm"
)

func GetAPICredentials(userID string) (*api.APICredentials, error) {

	var apiKeys models.ApiKeys

	result := config.DB.Where("user_id = ?", userID).First(&apiKeys)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.Log(slog.LevelError, "User not found", "UserID", userID)
			return nil, utils.CapitalizeError("user not found") // Or a more specific error
		}
		utils.Log(slog.LevelError, "Database error finding user", "Error", result.Error)
		return nil, utils.CapitalizeError("database error")
	}

	return &api.APICredentials{
		ClientSecret:    apiKeys.ClientSecret,
		ClientID:        apiKeys.ClientID,
		ClientSignature: apiKeys.OAuthSignature,
	}, nil
}
