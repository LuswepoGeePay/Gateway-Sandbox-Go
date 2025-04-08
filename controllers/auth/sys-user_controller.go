package auth

import (
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	adminservices "pg_sandbox/services/admin_services"
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

	users, err := adminservices.GetUsers(req)

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

	req := &auth.GetUsersRequest{
		Page:     int32(getRequest.Page),
		PageSize: int32(getRequest.PageSize),
	}

	users, err := adminservices.GetMerchants(req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("users"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Users"), gin.H{
		"users": users,
	})

}

func EditUserHandler(c *gin.Context) {
	var req auth.EditUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	err := adminservices.EditUser(&req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToUpdate("user"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessUpdate("User"))

}
