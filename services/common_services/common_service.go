package commonservices

import (
	tokenservices "pg_sandbox/services/token_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckEssentialHeaders(c *gin.Context) {

	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")
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

	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Content-Type is application/json"},
			},
		})
		c.Abort()
		return

	}

	if acceptedType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Accepted type is application/json"},
			},
		})
		c.Abort()
		return
	}

}
