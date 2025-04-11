package hostedcheckout

import (
	"pg_sandbox/proto/hcheckout"
	commonservices "pg_sandbox/services/common_services"
	hostedcheckoutservices "pg_sandbox/services/hosted_checkout_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func HostedCheckOutHandler(c *gin.Context) {

	xClientID := c.GetHeader("X-Client-Id")
	xTRef := c.GetHeader("X-Transaction-Ref")
	xCallbackUrl := c.GetHeader("X-Callback-URL")

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
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

	testCondition := c.Param("condition")

	if testCondition == "1" {

		utils.RespondWithSuccess(c, "Payment processed successfully")
		c.Abort()
		return

	}

	if testCondition == "2" {
		c.JSON(406, gin.H{
			"status":  "cancelled",
			"message": "Payment was cancelled",
		})
		c.Abort()
		return

	}

	if testCondition == "3" {
		utils.RespondWithError(c, 400, "Failed to process payment")
		c.Abort()
		return
	}

}
