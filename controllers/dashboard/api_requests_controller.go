package dashboard

import (
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GetAPIRequestsInfoHandler(c *gin.Context) {

	response, err := dashboardservices.GetApiStatistics()

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("API Statistics Card Information"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("API Statistics Card Info"), gin.H{
		"info": response,
	})

}

func GetAPIRequestsHandler(c *gin.Context) {

	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetAPIrequests{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	requests, err := dashboardservices.GetAPIRequests(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("requests"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("requests"), gin.H{
		"requests": requests,
	})
}
