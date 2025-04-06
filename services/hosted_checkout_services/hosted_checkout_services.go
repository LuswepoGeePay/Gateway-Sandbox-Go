package hostedcheckoutservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/hcheckout"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateCheckoutUrl(c *gin.Context, req *hcheckout.HCheckoutRequest, xClientID string, xTref string, xCallbackUrl string) {

	var existingClient models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientID).First(&existingClient)

	if result.Error != nil {
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

	generatedCheckoutUrl := models.CheckOutUrls{
		ID:            checkoutID,
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		CustomerName:  req.Customer.Name,
		CustomerEmail: req.Customer.Email,
		CancelUrl:     req.RedirectUrls.Cancel,
		SuccessUrl:    req.RedirectUrls.Success,
		FailedUrl:     req.RedirectUrls.Failure,
		GeneratedUrl:  checkoutUrl,
		TReference:    xTref,
	}

	result = tx.Create(&generatedCheckoutUrl)

	if result.Error != nil {
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

}
