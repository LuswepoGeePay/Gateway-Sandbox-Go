package collection

import (
	"pg_sandbox/proto/collection"
	collectionservices "pg_sandbox/services/collection_services"
	commonservices "pg_sandbox/services/common_services"

	"github.com/gin-gonic/gin"
)

func MakeCollectionHandler(c *gin.Context) {

	xClientId := c.GetHeader("X-Client-ID")
	xTransactionRef := c.GetHeader("X-Transaction-Ref")

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
		return
	}

	var req collection.CollectionRequest

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

	collectionservices.RequestToPay(c, xClientId, xTransactionRef, &req)
}
