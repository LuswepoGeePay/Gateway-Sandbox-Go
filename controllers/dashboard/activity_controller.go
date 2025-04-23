package dashboard

import (
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	dashboardservices "pg_sandbox/services/dashboard_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GetActivitiesHandler(c *gin.Context) {

	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &dashboard.GetActivityrequests{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	Activity, err := dashboardservices.GetActivity(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Activity"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Activity"), gin.H{
		"Activity": Activity,
	})
}
