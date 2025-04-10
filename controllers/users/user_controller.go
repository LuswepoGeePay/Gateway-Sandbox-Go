package users

import (
	"io/ioutil"
	"log"
	"net/http"
	"pg_sandbox/proto/api"
	"pg_sandbox/proto/user"
	commonservices "pg_sandbox/services/common_services"
	credentialsservices "pg_sandbox/services/credentials_services"
	userservices "pg_sandbox/services/user_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func NameLookUpHandler(c *gin.Context) {

	number := c.Param("number")

	commonservices.CheckEssentialHeaders(c)

	if c.IsAborted() {
		return
	}
	userservices.NameLookUp(c, number)

}

func SetFloatBalanceHander(c *gin.Context) {

	var req api.UpdateFloatReuest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	err := credentialsservices.SetFloatBalance(c, &req)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to set/update float balance", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Float Updated successfully")

}

func CallbackHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		c.String(http.StatusBadRequest, "can't read body")
		return
	}
	defer c.Request.Body.Close()

	log.Printf("Received callback: %s", string(body))
	c.String(http.StatusOK, "Callback received")
}

func ResetPasswordHandler(c *gin.Context) {
	var req user.ResetPassword

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	err = userservices.ResetPassword(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Failed to reset password.", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Password has been reset successfully")
}

func GetUserProfileHandler(c *gin.Context) {
	userId := c.Param("id")

	user, err := userservices.GetUserProfile(userId)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to fetch profile", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Retrieved user", gin.H{
		"user": user,
	})
}
