package auth

import (
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	authservices "pg_sandbox/services/auth_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GetUsersHandler(c *gin.Context) {
	var getRequest models.GetRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.RespondWithError(c, 400, utils.InvReqBody, err.Error())
		return
	}

	getRequest.SetDefaults()

	req := &auth.GetUsersRequest{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	users, err := authservices.GetUsers(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("users"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Users"), gin.H{
		"users": users,
	})

}
