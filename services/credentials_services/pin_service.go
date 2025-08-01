package credentialsservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/api"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func SetPin(req *api.SetPinRequest) error {

	if req.UserId == "" {
		return utils.CapitalizeError("user ID is required.")
	}
	if req.Pin == "" {
		return utils.CapitalizeError("pin is required.")
	}

	if len(req.Pin) < 8 {
		return utils.CapitalizeError("pin is supposed to be more than 8 digits")
	}

	tx := config.DB.Begin()

	result := tx.Where("user_id = ?", req.UserId).Model(&models.ApiKeys{}).Update("pin", req.Pin)

	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌Error", "unable to set pin merchant fee profile", "data", gin.H{
			"request_body": req,
		})
		return utils.CapitalizeError("unable to set new pin on merchant profile")
	}

	tx.Commit()

	return nil
}
