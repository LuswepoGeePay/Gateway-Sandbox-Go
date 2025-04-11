package disbursementservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"

	"github.com/gin-gonic/gin"
)

func CheckDisbursementBalance(c *gin.Context, xClientID string, xAuthSignature string) {

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

	var existingAuthSig models.ApiKeys

	result = config.DB.Where("o_auth_signature = ?", xAuthSignature).First(&existingAuthSig)

	if result.Error != nil {
		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Auth-Signature": []string{"The selected x-auth-signature is invalid."},
			},
		})
		return
	}

	var user models.User

	result = config.DB.Where("id = ? ", existingClient.UserID).First(&user)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"status":  "failed",
			"message": "User not",
			"error":   gin.H{"System error": []string{"Transaction Reference is invalid"}},
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Disbursement balance fetched successfully.",
		"data": gin.H{
			"merchant":     user.Fullname,
			"balance":      existingClient.FloatBalance,
			"currency":     "ZMW",
			"last_updated": existingClient.UpdatedAt,
		},
	})
}
