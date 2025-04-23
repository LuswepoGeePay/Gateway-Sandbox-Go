package auth

import (
	auth "pg_sandbox/proto/auth"
	authservices "pg_sandbox/services/auth_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {

	var req auth.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	userId, err := authservices.RegisterUser(c, &req)

	if err != nil {
		utils.RespondWithError(c, 400, utils.FailedToCreate("account"), err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessCreate("Account"), gin.H{
		"user_id": userId,
	})

}
