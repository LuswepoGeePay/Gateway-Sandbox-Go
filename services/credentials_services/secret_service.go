package credentialsservices

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/api"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GenerateSecret(req *api.RegenerateRequest) (*string, error) {

	secretBytes := make([]byte, 32)

	_, err := rand.Read(secretBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}
	secret := base64.StdEncoding.EncodeToString(secretBytes)

	tx := config.DB.Begin()

	var user models.User

	err = tx.Where("id = ?", req.UserId).Find(&user).Error
	if err != nil {
		utils.Log(slog.LevelError, "Unable to generate client secret, unable find user", "Error", err.Error)
		return nil, utils.CapitalizeError("cannot find user")
	}

	result := tx.Where("user_id = ?", user.ID).Model(&models.ApiKeys{}).Update("client_secret", secret)

	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌Error", "unable to add client secret to ", "data", gin.H{
			"error":        result.Error,
			"request_body": req,
		})
		return nil, utils.CapitalizeError("Error adding signature to profile")
	}

	tx.Commit()

	return &secret, nil
}
