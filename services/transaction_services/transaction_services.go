package transactionservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"

	"github.com/gin-gonic/gin"
)

func TransactionQuery(c *gin.Context, transationRef string) {

	var transaction models.Transactions

	if transationRef == "" {
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
		"message": "Transaction Status Retrieved",
		"data": gin.H{
			"status":   transaction.Status,
			"amount":   transaction.Amount,
			"customer": transaction.Customer,
			"channel":  transaction.Channel,
			"date":     transaction.Date,
		},
	})

}
