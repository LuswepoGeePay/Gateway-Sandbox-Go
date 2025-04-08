package auth

import (
	auth "pg_sandbox/proto/auth"
	authservices "pg_sandbox/services/auth_services"

	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	details, code, err := authservices.LoginUser(&req)

	if err != nil {
		utils.RespondWithError(c, code, "Failed to login", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Successful Login", gin.H{
		"user": details,
	})

}
