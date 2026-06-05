package dashboard

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"

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

	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetTransactionsRequest{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	userID, exists := c.Get("userID")
	merchantID := ""
	if exists {
		var user models.User
		if err := config.DB.Preload("Role").Where("id = ?", userID).First(&user).Error; err == nil {
			if user.Role.Name == "merchant" {
				merchantID = user.ID.String()
			}
		}
	}

	transactions, err := dashboardservices.GetTransactions(req, merchantID)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("transactions"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("transactions"), gin.H{
		"transactions": transactions,
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
