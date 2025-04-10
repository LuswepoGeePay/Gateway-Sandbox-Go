package disbursement

import (
	"pg_sandbox/proto/disbursement"
	commonservices "pg_sandbox/services/common_services"
	disbursementservices "pg_sandbox/services/disbursement_services"

	"github.com/gin-gonic/gin"
)

func MakeDisbursementHandler(c *gin.Context) {

	xClientID := c.GetHeader("X-Client-ID")
	xAuthSignature := c.GetHeader("X-Auth-Signature")
	xTRef := c.GetHeader("X-Transaction-Ref")
	xCallbackUrl := c.GetHeader("X-Callback-URL")

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
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

	disbursementservices.MakeDisbursement(c, &req, xClientID, xAuthSignature, xCallbackUrl, xTRef)

}

func QueryDisbursementHandler(c *gin.Context) {

	ref := c.Param("reference")

	xClientID := c.GetHeader("X-Client-ID")
	xAuthSignature := c.GetHeader("X-Auth-Signature")

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
		return
	}
	disbursementservices.QueryDisbursement(c, xClientID, xAuthSignature, ref)

}
