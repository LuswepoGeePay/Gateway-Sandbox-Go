package hostedcheckout

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
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

type CheckoutResponseRequest struct {
	CheckoutID  string `json:"checkout_id"`
	PhoneNumber string `json:"phone_number"`
	Amount      string `json:"amount"`
}

func HostedCheckoutResponseHandler(c *gin.Context) {
	var req CheckoutResponseRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.CheckoutID != "" {
		var checkout models.CheckOutUrls
		if err := config.DB.Where("id = ?", req.CheckoutID).First(&checkout).Error; err != nil {
			utils.RespondWithError(c, 404, "Checkout session not found")
			return
		}
	}

	network, err := utils.GetNetworkProvider(req.PhoneNumber)
	if err != nil {
		utils.RespondWithError(c, 400, "Invalid Phone number")
		return
	}

	var existingCheckoutRequest models.CheckOutUrls
	if err := config.DB.Where("id = ?", req.CheckoutID).First(&existingCheckoutRequest).Error; err != nil {
		utils.RespondWithError(c, 404, "Checkout session not found")
		return
	}

	var existingTransaction models.Transactions
	if err := config.DB.Where("reference = ?", existingCheckoutRequest.TransactionReference).First(&existingTransaction).Error; err != nil {
		utils.RespondWithError(c, 404, "Transaction not found")
		return
	}

	updates := map[string]interface{}{
		"customer": req.PhoneNumber,
		"amount":   req.Amount,
		"status":   "successful",
		"channel":  network,
	}

	if err := config.DB.Model(&models.Transactions{}).Where("id = ?", existingTransaction.ID).Updates(updates).Error; err != nil {
		utils.RespondWithError(c, 500, "Failed to update transaction status")
		return
	}

	tenDigitRand := utils.GenerateTenDigitCode()
	c.JSON(200, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Payment was successful.",
		"data": gin.H{
			"status":                "successful",
			"transaction_reference": existingTransaction.Reference,
			"external_reference":    tenDigitRand,
			"return_url":            existingCheckoutRequest.ReturnUrl,
		},
	})

}

func GetCheckoutSession(c *gin.Context) {

	userId := c.Query("user_id")

	var checkouts []models.CheckOutUrls
	if err := config.DB.Where("user_id = ?", userId).Find(&checkouts).Error; err != nil {
		utils.RespondWithError(c, 404, "Checkouts not found")
		return
	}

	utils.RespondWithSuccess(c, "checkouts fetched", gin.H{
		"checkouts": checkouts,
	})

}
