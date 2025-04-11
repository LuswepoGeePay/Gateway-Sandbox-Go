package disbursementservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/disbursement"
	"pg_sandbox/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MakeDisbursement(c *gin.Context, req *disbursement.DisbursementRequest, xClientID string, xAuthSignature string, xCallbackUrl string, xTref string) {

	var existingClient models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClient)

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

	if req.Amount == "" {

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

	amountConverted, err := strconv.ParseFloat(req.Amount, 64)
	tCode := utils.GenerateTenDigitCode()

	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"amount": []string{"Invalid amount format."},
			},
		})
		return
	}

	if amountConverted <= 0 {
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

	if amountConverted > floatBalancedConverted {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Insufficient disbursement balance for transaction, please request for more float",
		})
		return
	}

	tx := config.DB.Begin()

	newDisbursementBalance := floatBalancedConverted - amountConverted

	result = tx.Model(&models.ApiKeys{}).Where("user_id = ?", existingClient.UserID).Update("float_balance", strconv.FormatFloat(newDisbursementBalance, 'f', -1, 64))

	if result.Error != nil {

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

	transaction := models.Transactions{
		ID:        uuid.New(),
		Reference: xTref,
		Channel:   network,
		Customer:  req.PhoneNumber,
		Amount:    req.Amount,
		Status:    "completed",
		Narration: req.Narration,
		Type:      "disbursement",
	}

	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()
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

	if xCallbackUrl != "" {
		CallbackHandler(xCallbackUrl, models.CallbackPayload{
			Code:          200,
			Status:        "successful",
			Message:       "The shackles have been sent man",
			TransactionID: xTref,
			ExternalID:    tCode,
		})
	}

	tx.Commit()

	if network == "mtn" || network == "airtel" {
		c.JSON(200, gin.H{
			"code":    200,
			"status":  "successful",
			"message": "Disbursement was processed successfully",
			"data": gin.H{
				"transaction_id": xTref,
				"external_id":    tCode,
			},
		})
		return
	}

	if network == "zamtel" {
		c.JSON(200, gin.H{
			"code":    200,
			"status":  "successful",
			"message": "Disbursement was processed successfully",
			"data": gin.H{
				"transaction_id": xTref,
				"external_id":    tCode,
			},
		})
		return
	}

}
