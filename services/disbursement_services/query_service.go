package disbursementservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func QueryDisbursement(c *gin.Context, xClientID string, xAuthSig string, Tref string) {

	var existingClientID models.ApiKeys
	start := time.Now()

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClientID)

	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Client ID is invalid")

		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Client-ID": []string{"The selected x-client-id is invalid."},
			},
		})
		return
	}

	var existingAuthSig models.ApiKeys

	result = config.DB.Where("o_auth_signature = ?", xAuthSig).First(&existingAuthSig)

	if result.Error != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/disburse/status/", "GET", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Auth-Signature": []string{"The selected x-auth-signature is invalid."},
			},
		})
		return
	}

	var transaction models.Transactions

	if Tref == "" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/disburse/status/", "GET", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "failed",
			"message": "Invalid Transaction Reference",
			"error":   []string{"Transaction Reference is invalid"},
		})
		return

	}

	result = config.DB.Where("reference = ? AND type = ?", Tref, "disbursement").First(&transaction)

	if result.Error != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/disburse/status/", "GET", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(404, gin.H{
			"code":    404,
			"status":  "failed",
			"message": "Transaction Not Found",
			"error":   gin.H{"Transaction Reference": []string{"Transaction Reference is invalid"}},
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Disbursement status fetched successfully.",
		"data": gin.H{
			"status":                transaction.Status,
			"transaction_reference": transaction.Reference,
		},
	})

	elapsed := time.Since(start).Milliseconds()
	logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/disburse/status/", "GET", "failed", strconv.FormatInt(elapsed, 10))

}
