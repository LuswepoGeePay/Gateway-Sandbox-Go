package collectionservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/collection"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestToPay(c *gin.Context, xClientId string, xTransactionRef string, req *collection.CollectionRequest) {

	start := time.Now()
	// handle logic

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientId).First(&existingClientID)

	if result.Error != nil {

		utils.Log(slog.LevelError, "Error", "Client ID is invalid")
		c.JSON(402, gin.H{
			"code":    402,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Client-ID": []string{"The selected x-client-id is invalid."},
			},
		})
		return
	}

	if strings.TrimSpace(xTransactionRef) == "" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(402, gin.H{
			"code":    402,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Transaction-Ref": []string{"Transaction reference cannot be empty."},
			},
		})
		return
	}

	var existingTransaction models.Transactions

	result = config.DB.Where("reference = ?", xTransactionRef).First(&existingTransaction)

	if result.Error == nil { // Only return an error if the transaction exists
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Transaction-Ref": []string{"The x-transaction-ref has already been taken."},
			},
		})
		return
	}

	network, err := utils.GetNetworkProvider(req.PhoneNumber)

	tCode := utils.GenerateTenDigitCode()

	if err != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"phone_number": []string{"Invalid Phone number"},
			},
		})
		return

	}

	tStatus := ""

	if network == "mtn" {
		tStatus = "pending"
	} else if network == "airtel" {
		tStatus = "pending"
	} else if network == "zamtel" {
		tStatus = "sucessful"
	}

	transaction := models.Transactions{
		ID:        uuid.New(),
		Reference: xTransactionRef,
		Channel:   network,
		Customer:  req.PhoneNumber,
		Amount:    req.Amount,
		Status:    tStatus,
		Type:      "collection",
		Date:      time.Now(),
	}

	tx := config.DB.Begin()

	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "System Error",
			"errors": gin.H{
				"Tranaction Error": []string{"Failed to create transaction"},
			},
		})
		return
	}

	tx.Commit()

	if network == "mtn" || network == "airtel" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "success", strconv.FormatInt(elapsed, 10))

		c.JSON(202, gin.H{
			"code":    202,
			"status":  "pending",
			"message": "Request sent. Awaiting customer action.",
			"data": gin.H{
				"transaction_reference": xTransactionRef,
				"external_reference":    tCode,
			},
		})
		return
	}

	if network == "zamtel" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "success", strconv.FormatInt(elapsed, 10))

		c.JSON(200, gin.H{
			"code":    200,
			"status":  "success",
			"message": "Payment was successful",
			"data": gin.H{
				"transaction_reference": xTransactionRef,
				"external_reference":    tCode,
			},
		})
		return
	}

}
