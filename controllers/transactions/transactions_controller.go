package transactions

import (
	tokenservices "pg_sandbox/services/token_services"
	transactionservices "pg_sandbox/services/transaction_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func TransactionQueryHandler(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(401, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	err := tokenservices.ValidateOAuthToken(tokenString)
	if err != nil {
		utils.RespondWithError(c, 401, "Invalid Token")
		c.Abort() // Abort the request to prevent further execution
		return
	}

	transactionID := c.Param("id")

	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")

	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Content-Type is application/json"},
			},
		})

		return
	}

	if acceptedType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Accept is application/json"},
			},
		})

		return
	}

	transactionservices.TransactionQuery(c, transactionID)

}
