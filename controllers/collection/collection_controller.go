package collection

import (
	"pg_sandbox/proto/collection"
	collectionservices "pg_sandbox/services/collection_services"
	tokenservices "pg_sandbox/services/token_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func MakeCollectionHandler(c *gin.Context) {

	xClientId := c.GetHeader("X-Client-ID")
	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")
	xTransactionRef := c.GetHeader("X-Transaction-Ref")
	authorization := c.GetHeader("Authorization")

	if authorization == "" {
		c.JSON(400, gin.H{
			"message": "unauthenticated",
		})

		return
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	err := tokenservices.ValidateOAuthToken(tokenString)
	if err != nil {
		utils.RespondWithError(c, 401, "Invalid Token")
		c.Abort()
		return
	}

	var req collection.CollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid Request Body.",
			"errors": gin.H{
				"Incomplete Body": []string{"The JSON body is missing"},
			},
		})

		return

	}

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

	collectionservices.RequestToPay(c, xClientId, xTransactionRef, &req)
}
