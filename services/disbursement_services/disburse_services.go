package disbursementservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/disbursement"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MakeDisbursement(c *gin.Context, req *disbursement.DisbursementRequest, xClientID string, xAuthSignature string, xCallbackUrl string, xTref string) {

	var existingClient models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClient)
	start := time.Now()

	if result.Error != nil {

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

	result = config.DB.Where("o_auth_signature = ?", xAuthSignature).First(&existingAuthSig)

	if result.Error != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

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

	if strings.TrimSpace(xTref) == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Transaction-Ref": []string{"Transaction reference cannot be empty."},
			},
		})
		return
	}

	var existingTransaction models.Transactions

	result = config.DB.Where("reference = ?", xTref).First(&existingTransaction)

	if result.Error == nil { // Only return an error if the transaction exists
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

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

	// tCode := utils.GenerateTenDigitCode()

	if err != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

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

	if req.Amount == 0 {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"amount": []string{"Amount is required."},
			},
		})
		return

	}

	if req.Amount <= 0 {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"amount": []string{"Invalid amount. Amount must be greater than 0"},
			},
		})
		return
	}

	floatBalancedConverted, err := strconv.ParseFloat(existingClient.FloatBalance, 64)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"code":    "500",
			"error":   "Error detecting balance",
			"message": "Unable to convert float balance",
		})

		return

	}

	if float64(req.Amount) > floatBalancedConverted {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Insufficient disbursement balance, please top up",
		})
		return
	}

	tx := config.DB.Begin()

	newDisbursementBalance := floatBalancedConverted - float64(req.Amount)

	result = tx.Model(&models.ApiKeys{}).Where("user_id = ?", existingClient.UserID).Update("float_balance", strconv.FormatFloat(newDisbursementBalance, 'f', -1, 64))

	if result.Error != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(500, gin.H{
			"code":    500,
			"status":  "error",
			"message": "Server error.",
			"errors": gin.H{
				"Float": []string{"Unable to update float balance."},
			},
		})
		return
	}

	tStatus := ""

	// if req.IsFailed {
	// 	tStatus = "failed"
	// }

	tStatus = "successful"

	transaction := models.Transactions{
		ID:        uuid.New(),
		Reference: xTref,
		Channel:   network,
		Customer:  req.PhoneNumber,
		Amount:    string(req.Amount),
		Status:    tStatus,
		Narration: req.Narration,
		Type:      "disbursement",
		Date:      time.Now(),
		UserID:    existingClient.UserID,
	}

	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(500, gin.H{
			"code":    500,
			"status":  "error",
			"message": "System Error",
			"errors": gin.H{
				"Tranaction Error": []string{"Failed to create transaction"},
			},
		})
		return
	}

	tCode := utils.GenerateTenDigitCode()

	if xCallbackUrl != "" {
		CallbackHandler(xCallbackUrl, models.CallbackPayload{
			Code:    200,
			Status:  "successful",
			Message: "Disbursement has been successfully processed and settled.",
			Data: models.CallbackPayloadData{
				TransactionReference: xTref,
				ExternalReference:    tCode,
				Customer:             req.PhoneNumber,
				Amount:               string(req.Amount),
			},
		})
	}

	// if xCallbackUrl != "" && req.IsFailed {
	// 	CallbackHandler(xCallbackUrl, models.CallbackPayload{
	// 		Code:    402,
	// 		Status:  "failed",
	// 		Message: "Transaction failed. Please try again.",
	// 		Data: models.CallbackPayloadData{
	// 			TransactionReference: xTref,
	// 			ExternalReference:    "",
	// 			Customer:             req.PhoneNumber,
	// 			Amount:               string(req.Amount),
	// 		},
	// 	})
	// }

	tx.Commit()

	if network == "mtn" || network == "airtel" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "success", strconv.FormatInt(elapsed, 10))

		c.JSON(200, gin.H{
			"code":    200,
			"status":  "successful",
			"message": "Disbursement has been successfully processed and settled.",
			"data": gin.H{
				"transaction_id":     xTref,
				"external_reference": tCode,
			},
		})
		return
	}

	if network == "zamtel" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/mobile-money/disburse", "POST", "success", strconv.FormatInt(elapsed, 10))

		c.JSON(200, gin.H{
			"code":    200,
			"status":  "successful",
			"message": "Disbursement has been successfully processed and settled.",
			"data": gin.H{
				"transaction_id": xTref,
				"external_id":    tCode,
			},
		})
		return
	}

}
