package credentialsservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/api"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func SetFloatBalance(c *gin.Context, req *api.UpdateFloatReuest) error {

	if req.UserId == "" {
		utils.Log(slog.LevelError, "❌Error", "Unable to set float balance, user id is required")
		return utils.CapitalizeError("user ID is required.")
	}
	if req.Float == "" {
		utils.Log(slog.LevelError, "❌Error", "Unable to set float balance float amount is required")
		return utils.CapitalizeError("float is required.")
	}

	tx := config.DB.Begin()

	// var user models.ApiKey

	result := tx.Where("user_id = ?", req.UserId).Model(&models.ApiKeys{}).Update("float_balance", req.Float)

	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌Error", "Unable to update float balance")
		return utils.CapitalizeError("unable to set new float balance on merchant profile")
	}
	tx.Commit()
	return nil
}
