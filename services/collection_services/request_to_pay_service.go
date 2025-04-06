package collectionservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/collection"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestToPay(c *gin.Context, xClientId string, xTransactionRef string, req *collection.CollectionRequest) {

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientId).First(&existingClientID)

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

	if strings.TrimSpace(xTransactionRef) == "" {
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

	result = config.DB.Where("reference = ?", xTransactionRef).First(&existingTransaction)

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

	tCode := utils.GenerateTenDigitCode()

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
		c.JSON(202, gin.H{
			"code":    202,
			"status":  "pending",
			"message": "Request sent. Awaiting customer action.",
			"data": gin.H{
				"transaction_id": xTransactionRef,
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
				"transaction_id": xTransactionRef,
				"external_id":    tCode,
			},
		})
		return
	}

}
