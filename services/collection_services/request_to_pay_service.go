package collectionservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/collection"
	disbursementservices "pg_sandbox/services/disbursement_services"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestToPay(c *gin.Context, xClientId string, xTransactionRef string, xCallbackUrl string, req *collection.CollectionRequest) {

	start := time.Now()
	// handle logic

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientId).First(&existingClientID)

	if result.Error != nil {

		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, client ID is invalid", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})
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

		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, transaction reference is empty", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})

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

		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, transaction reference has already been taken", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})

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
		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, inavlid phone number", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})
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

	// if req.IsFailed {

	// 	tStatus = "failed"
	// }

	if network == "mtn" {
		tStatus = "pending"
	} else if network == "airtel" {
		tStatus = "pending"
	} else if network == "zamtel" {
		tStatus = "successful"
	}

	transaction := models.Transactions{
		ID:        uuid.New(),
		Reference: xTransactionRef,
		Channel:   network,
		Customer:  req.PhoneNumber,
		Amount:    string(req.Amount),
		Status:    tStatus,
		Type:      "collection",
		Date:      time.Now(),
		UserID:    existingClientID.UserID,
	}

	tx := config.DB.Begin()

	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))
		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, unable to create transaction", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
			"error":           err,
		})
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

	if xCallbackUrl != "" && tStatus == "failed" {
		disbursementservices.CallbackHandler(xCallbackUrl, models.CallbackPayload{
			Code:    200,
			Status:  "failed",
			Message: "Transaction failed. Please try again.",
			Data: models.CallbackPayloadData{
				TransactionReference: xTransactionRef,
				ExternalReference:    "",
				Customer:             req.PhoneNumber,
				Amount:               string(req.Amount),
			},
		})
	}

	if tStatus == "pending" || tStatus == "successful" {

		if xCallbackUrl != "" {
			disbursementservices.CallbackHandler(xCallbackUrl, models.CallbackPayload{
				Code:    200,
				Status:  "successful",
				Message: "Transaction has been successfully processed and settled.",
				Data: models.CallbackPayloadData{
					TransactionReference: xTransactionRef,
					ExternalReference:    tCode,
					Customer:             req.PhoneNumber,
					Amount:               string(req.Amount),
				},
			})
		}
	}

	// if !req.IsFailed {
	// 	elapsed := time.Since(start).Milliseconds()
	// 	logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "failed", strconv.FormatInt(elapsed, 10))

	// 	c.JSON(200, gin.H{
	// 		"code":    402,
	// 		"status":  "failed",
	// 		"message": "Payment Failed: low balance or payee limit reached or not allowed",
	// 		"data": gin.H{
	// 			"transaction_reference": xTransactionRef,
	// 		},
	// 	})
	// 	return
	// }

	if network == "mtn" || network == "airtel" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/mobile-money/collect", "POST", "success", strconv.FormatInt(elapsed, 10))

		c.JSON(202, gin.H{
			"code":    200,
			"status":  "pending",
			"message": "Request sent. Awaiting customer action.",
			"data": gin.H{
				"transaction_reference": xTransactionRef,
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
