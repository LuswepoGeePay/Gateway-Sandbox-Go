package dashboard

import (
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GetUserStatisticsHandler(c *gin.Context) {

	response, err := dashboardservices.GetUserStatistics()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("User statistics Card Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("User statistics Card Info"), gin.H{
		"info": response,
	})

}

func GetMerchantStatisticsHandler(c *gin.Context) {

	response, err := dashboardservices.GetMerchantStatistics()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Merchant statistics Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Merchant statistics Card Info"), gin.H{
		"info": response,
	})

}
