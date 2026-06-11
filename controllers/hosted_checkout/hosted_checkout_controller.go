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
	// var req CheckoutResponseRequest
	// if err := c.ShouldBindJSON(&req); err == nil && req.CheckoutID != "" {
	// 	var checkout models.CheckOutUrls
	// 	if err := config.DB.Where("id = ?", req.CheckoutID).First(&checkout).Error; err != nil {
	// 		utils.RespondWithError(c, 404, "Checkout session not found")
	// 		return
	// 	}

	// 	network, err := utils.GetNetworkProvider(req.PhoneNumber)
	// 	if err != nil {
	// 		utils.RespondWithError(c, 400, "Invalid Phone number")
	// 		return
	// 	}

	// 	TReference := uuid.New()

	// 	var existingTx models.Transactions
	// 	if err := config.DB.Where("reference = ?", TReference).First(&existingTx).Error; err == nil {
	// 		utils.RespondWithError(c, 400, "Transaction has already been initiated/processed")
	// 		return
	// 	}

	// 	txID := uuid.New()
	// 	tStatus := "successful"

	// 	transaction := models.Transactions{
	// 		ID:          txID,
	// 		Reference:   TReference.String(),
	// 		Channel:     network,
	// 		Customer:    req.PhoneNumber,
	// 		Amount:      checkout.Amount,
	// 		Status:      tStatus,
	// 		Type:        "collection",
	// 		Date:        time.Now(),
	// 		UserID:      checkout.UserID,
	// 		CallbackUrl: checkout.CallbackUrl,
	// 	}

	// 	if err := config.DB.Create(&transaction).Error; err != nil {
	// 		utils.RespondWithError(c, 500, "Failed to create transaction")
	// 		return
	// 	}

	// 	c.JSON(200, gin.H{
	// 		"code":    200,
	// 		"status":  "success",
	// 		"message": "Transaction processed successfully. Redirecting to Merchant Page",
	// 		"data": gin.H{
	// 			"transaction_reference": TReference.String(),
	// 			"transaction_id":        txID.String(),
	// 			"status":                tStatus,
	// 			"redirect_url":          checkout.ReturnUrl,
	// 		},
	// 	})
	// 	return
	// }

	// testCondition := c.Param("condition")

	// if testCondition == "1" {
	// 	utils.RespondWithSuccess(c, "Payment processed successfully")
	// 	c.Abort()
	// 	return
	// }

	// if testCondition == "2" {
	// 	c.JSON(406, gin.H{
	// 		"status":  "cancelled",
	// 		"message": "Payment was cancelled",
	// 	})
	// 	c.Abort()
	// 	return
	// }

	// if testCondition == "3" {
	// 	utils.RespondWithError(c, 400, "Failed to process payment")
	// 	c.Abort()
	// 	return
	// }

	var req CheckoutResponseRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.CheckoutID != "" {
		var checkout models.CheckOutUrls
		if err := config.DB.Where("id = ?", req.CheckoutID).First(&checkout).Error; err != nil {
			utils.RespondWithError(c, 404, "Checkout session not found")
			return
		}
	}

	// network, err := utils.GetNetworkProvider(req.PhoneNumber)
	// if err != nil {
	// 	utils.RespondWithError(c, 400, "Invalid Phone number")
	// 	return
	// }

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
