package users

import (
	"io/ioutil"
	"log"
	"net/http"
	"pg_sandbox/proto/api"
	"pg_sandbox/proto/user"
	credentialsservices "pg_sandbox/services/credentials_services"
	tokenservices "pg_sandbox/services/token_services"
	userservices "pg_sandbox/services/user_services"
	"pg_sandbox/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func NameLookUpHandler(c *gin.Context) {

	number := c.Param("number")
	authorization := c.GetHeader("Authorization")

	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")

	if authorization == "" {
		c.JSON(400, gin.H{
			"message": "unauthenticated",
		})

		return
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	err := tokenservices.ValidateOAuthToken(tokenString)
	if err != nil {
		utils.RespondWithError(c, 401, "Invalid Token")
		c.Abort()
		return
	}

	if contentType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Content-Type is application/json"},
			},
		})

		return
	}

	if acceptedType != "application/json" {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Accept is application/json"},
			},
		})

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
