package disbursement

import (
	"pg_sandbox/proto/disbursement"
	disbursementservices "pg_sandbox/services/disbursement_services"
	tokenservices "pg_sandbox/services/token_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func MakeDisbursementHandler(c *gin.Context) {

	xClientID := c.GetHeader("X-Client-ID")
	xAuthSignature := c.GetHeader("X-Auth-Signature")
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

	var req disbursement.DisbursementRequest

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

	disbursementservices.MakeDisbursement(c, &req, xClientID, xAuthSignature, xCallbackUrl, xTRef)

}

func QueryDisbursementHandler(c *gin.Context) {

	ref := c.Param("reference")

	xClientID := c.GetHeader("X-Client-ID")
	xAuthSignature := c.GetHeader("X-Auth-Signature")
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

	disbursementservices.QueryDisbursement(c, xClientID, xAuthSignature, ref)

}
