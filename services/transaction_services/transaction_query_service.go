package transactionservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func TransactionQuery(c *gin.Context, transationRef string) {

	var transaction models.Transactions

	if transationRef == "" {
		utils.Log(slog.LevelError, "❌Error", "Invalid transaction reference", "endpoint", "/v1/mobile-money/check-status/", "reference", transationRef)

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "failed",
			"message": "Invalid Transaction Reference",
			"error":   []string{"Transaction Reference is invalid"},
		})
		return

	}

	result := config.DB.Where("reference = ? AND type = ?", transationRef, "collection").First(&transaction)

	if result.Error != nil {
		utils.Log(slog.LevelError, "❌Error", "Invalid transaction reference", "endpoint", "/v1/mobile-money/check-status/", "reference", transationRef)

		c.JSON(404, gin.H{
			"code":    404,
			"status":  "error",
			"message": "Transaction Not Found",
		})
		return
	}

	tCode := utils.GenerateTenDigitCode()

	c.JSON(200, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Transaction Status Fetched successfully",
		"data": gin.H{
			"status":                "successful",
			"message":               "Transaction was processed successfully",
			"transaction_reference": transaction.Reference,
			"external_reference":    tCode,
		},
	})

	utils.Log(slog.LevelInfo, "✅Info", "Transaction Status Retrieved", "endpoint", "/v1/mobile-money/check-status/", "reference", transationRef)

}
