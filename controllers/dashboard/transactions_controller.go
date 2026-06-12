package dashboard

import (
	"pg_sandbox/proto/dashboard"
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTransactionInfoHandler(c *gin.Context) {

	response, err := dashboardservices.GetTransactionStatistics()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Transaction Card Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Transaction Card Info"), gin.H{
		"info": response,
	})

}

func GetTransactionsHandler(c *gin.Context) {

	status := c.Query("status")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	transactionReference := c.Query("transaction_reference")
	customer := c.Query("customer")
	externalReference := c.Query("external_reference")
	transactionType := c.Query("transaction_type")
	channel := c.Query("channel")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	req := dashboard.GetTransactionsRequest{
		Page:                 int32(pageInt),
		PageSize:             int32(pageSizeInt),
		Status:               status,
		TransactionReference: transactionReference,
		Customer:             customer,
		ExternalReference:    externalReference,
		TransactionType:      transactionType,
		Channel:              channel,
		StartDate:            startDate,
		EndDate:              endDate,
	}

	transactions, err := dashboardservices.GetTransactions(&req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("transactions"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("transactions"), gin.H{
		"data": transactions,
	})
}

func GetTransactionsChannelHandler(c *gin.Context) {

	response, err := dashboardservices.GetTransactionChannelStats()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Transaction Channel Statistics Card Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Transaction Channel Statistics Card Info"), gin.H{
		"info": response,
	})

}
