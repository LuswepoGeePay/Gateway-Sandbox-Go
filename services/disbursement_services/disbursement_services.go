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

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClientID)

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

	if req.Narration == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"narration": []string{"Invalid narration. Please provide a narration."},
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
		Status:    "successful",
		Narration: req.Narration,
		Type:      "disbursement",
	}

	tx := config.DB.Begin()

	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()

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
		c.JSON(200, gin.H{
			"code":    200,
			"status":  "successful",
			"message": "Funds disbursed",
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
			"message": "Process service request successfully.",
			"data": gin.H{
				"transaction_id": xTref,
				"external_id":    tCode,
			},
		})
		return
	}

}

func QueryDisbursement(c *gin.Context, xClientID string, xAuthSig string, Tref string) {

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClientID)

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

	result = config.DB.Where("o_auth_signature = ?", xAuthSig).First(&existingAuthSig)

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

	var transaction models.Transactions

	if Tref == "" {
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
		"message": "Disbursement Status Retrieved",
		"data": gin.H{
			"status":    transaction.Status,
			"amount":    transaction.Amount,
			"customer":  transaction.Customer,
			"channel":   transaction.Channel,
			"date":      transaction.Date,
			"narration": transaction.Narration,
		},
	})

}
