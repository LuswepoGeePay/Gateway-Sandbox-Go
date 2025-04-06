package hostedcheckout

import (
	"pg_sandbox/proto/hcheckout"
	hostedcheckoutservices "pg_sandbox/services/hosted_checkout_services"
	tokenservices "pg_sandbox/services/token_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func HostedCheckOutHandler(c *gin.Context) {

	xClientID := c.GetHeader("X-Client-Id")
	xTRef := c.GetHeader("X-Transaction-Ref")
	xCallbackUrl := c.GetHeader("X-Callback-URL")
	authorization := c.GetHeader("Authorization")
	acceptedH := c.GetHeader("Accept")
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

	var req hcheckout.HCheckoutRequest
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

	if acceptedH != "application/json" {
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

	hostedcheckoutservices.GenerateCheckoutUrl(c, &req, xClientID, xTRef, xCallbackUrl)

}

func GetHostedCheckoutDetailsHandler(c *gin.Context) {

	id := c.Param("id")

	response, err := hostedcheckoutservices.GetCheckoutDetails(id)

	if err != nil {
		utils.RespondWithError(c, 400, err.Error())
	}

	utils.RespondWithSuccess(c, "checkout details fetched", gin.H{
		"response": response,
	})

}

func HostedCheckoutResponseHandler(c *gin.Context) {

}
