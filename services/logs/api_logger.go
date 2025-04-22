package logs

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LogApiCall(c *gin.Context, userID string, endpoint string, method string, status string, responseTimeMs string) error {

	userId, err := uuid.Parse(userID)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "Failed to parse user ID", endpoint, "unable to log API call")
		return utils.CapitalizeError(err.Error())
	}

	apiLog := models.APILogs{
		ID:           uuid.New(),
		UserID:       userId,
		Endpoint:     endpoint,
		Method:       method,
		Status:       status,
		ResponseTime: responseTimeMs,
		IPAddress:    utils.GetIPAddress(c),
	}

	err = config.DB.Create(&apiLog).Error
	if err != nil {
		utils.Log(slog.LevelError, "Failed to log API Call:", err.Error())
		return err
	}

	return nil

}

func LogActivity(c *gin.Context, userID string, action string, entityType string, entityID string, requestBody string, responseCode int, status string, errorMsg string, responseTimeMs string) error {

	userId, err := uuid.Parse(userID)

	if err != nil {
		utils.Log(slog.LevelError, "Error", "Failed to parse user ID", userID, "unable to log Activity call")
		return utils.CapitalizeError(err.Error())
	}

	apiLog := models.ActivityLogs{
		ID:        uuid.New(),
		UserID:    userId,
		Action:    action,
		Entity:    entityType,
		EntityID:  entityID,
		IPAddress: utils.GetIPAddress(c),
	}

	err = config.DB.Create(&apiLog).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error: Failed to log Activity", err.Error())
		return err
	}

	return nil

}
