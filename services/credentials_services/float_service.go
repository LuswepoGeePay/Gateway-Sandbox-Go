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
		return utils.CapitalizeError("user ID is required.")
	}
	if req.Float == "" {
		return utils.CapitalizeError("float is required.")
	}

	tx := config.DB.Begin()

	// var user models.ApiKey

	result := tx.Where("user_id = ?", req.UserId).Model(&models.ApiKeys{}).Update("float_balance", req.Float)

	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "Error adding to db", "Error", result.Error)
		return utils.CapitalizeError("unable to set new float balance on merchant profile")
	}
	tx.Commit()
	return nil
}
