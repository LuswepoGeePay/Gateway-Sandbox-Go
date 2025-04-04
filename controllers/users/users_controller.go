package users

import (
	tokenservices "pg_sandbox/services/token_services"
	userservices "pg_sandbox/services/user_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func NameLookUpHandler(c *gin.Context) {

	number := c.Param("number")
	authorization := c.GetHeader("Authorization")

	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")

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

	userservices.NameLookUp(c, number)

}
