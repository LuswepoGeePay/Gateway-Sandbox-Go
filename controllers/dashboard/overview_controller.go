package dashboard

import (
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GetOverviewCardsInfoHandler(c *gin.Context) {

	response, err := dashboardservices.GetOverviewCardsInfo()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Overview Card Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Overview Card Info"), gin.H{
		"info": response,
	})

}

func GetUsersHandler(c *gin.Context) {
	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetUsersRequest{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	users, err := dashboardservices.GetUsers(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("users"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Users"), gin.H{
		"users": users,
	})

}

func GetMerchantsHandler(c *gin.Context) {
	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetUsersRequest{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	users, err := dashboardservices.GetMerchants(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("users"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Users"), gin.H{
		"users": users,
	})

}

func GetTopMerchantsHandler(c *gin.Context) {
	response, err := dashboardservices.GetTopMerchants()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Top Merchant Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Top Merchant Info"), gin.H{
		"info": response,
	})

}
