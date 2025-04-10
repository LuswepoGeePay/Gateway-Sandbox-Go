package transactions

import (
	commonservices "pg_sandbox/services/common_services"
	transactionservices "pg_sandbox/services/transaction_services"

	"github.com/gin-gonic/gin"
)

func TransactionQueryHandler(c *gin.Context) {

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
		return
	}

	transactionID := c.Param("id")
	transactionservices.TransactionQuery(c, transactionID)

}
