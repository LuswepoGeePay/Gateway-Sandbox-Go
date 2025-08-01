package hostedcheckoutservices

import (
	"log/slog"
	"net/url"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/hcheckout"
	disbursementservices "pg_sandbox/services/disbursement_services"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateCheckoutUrl(c *gin.Context, req *hcheckout.HCheckoutRequest, xClientID string, xTref string, xCallbackUrl string) {

	var existingClient models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClient)

	start := time.Now()

	if result.Error != nil {
		utils.Log(slog.LevelError, "Error", "Client ID is invalid")
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

	if strings.TrimSpace(xTref) == "" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/checkout/session", "POST", "failed", strconv.FormatInt(elapsed, 10))

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

	result = config.DB.Where("reference = ?", xTref).First(&existingTransaction)

	if result.Error == nil { // Only return an error if the transaction exists
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/checkout/session", "POST", "failed", strconv.FormatInt(elapsed, 10))

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

	tx := config.DB.Begin()

	checkoutID := uuid.New()

	checkoutUrl := req.CheckoutBaseUrl + checkoutID.String()

	returnUrl := ""

	var returnUrlParsed *url.URL
	var err error

	if req.ReceiptRedirect {
		returnUrlParsed, err = url.Parse(req.ReturnUrl)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"status":  "error",
				"message": "Invalid return URL provided",
			})
			return
		}

		txCode := utils.GenerateTenDigitCode()
		params := returnUrlParsed.Query()
		params.Set("status", "successful")
		params.Set("message", "Your transaction was completed successfully.") // customize if needed
		params.Set("transaction_reference", xTref)
		params.Set("external_reference", txCode)

		returnUrlParsed.RawQuery = params.Encode()
	}

	if req.ReceiptRedirect {
		returnUrl = returnUrlParsed.String()
	} else {
		returnUrl = req.ReturnUrl
	}

	generatedCheckoutUrl := models.CheckOutUrls{
		ID:            checkoutID,
		OrderID:       req.OrderId,
		Amount:        strconv.FormatFloat(float64(req.Amount), 'f', -1, 64),
		CustomerName:  req.Customer.Name,
		CustomerEmail: req.Customer.Email,
		ReturnUrl:     returnUrl,
		GeneratedUrl:  checkoutUrl,
		TReference:    xTref,
	}

	result = tx.Create(&generatedCheckoutUrl)

	if result.Error != nil {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClient.UserID.String(), "/v1/checkout/session", "POST", "failed", strconv.FormatInt(elapsed, 10))

		c.JSON(500, gin.H{
			"code":    500,
			"status":  "error",
			"message": "Server error.",
			"errors": gin.H{
				"URL": []string{"Unable to save url."},
			},
		})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{
		"status":       "success",
		"message":      "Checkout session created",
		"checkout_url": checkoutUrl,
	})

	if xCallbackUrl != "" {
		disbursementservices.CallbackHandler(xCallbackUrl, models.CallbackPayload{
			Code:    200,
			Status:  "successful",
			Message: "Transaction successful. Please try again.",
			Data: models.CallbackPayloadData{
				TransactionReference: xTref,
				ExternalReference:    "",
				Customer:             "",
				Amount:               strconv.FormatFloat(float64(req.Amount), 'f', -1, 32),
			},
		})
	}

	elapsed := time.Since(start).Milliseconds()
	logs.LogApiCall(c, existingClient.UserID.String(), "/v1/checkout/session", "POST", "success", strconv.FormatInt(elapsed, 10))

}
