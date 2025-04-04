package auth

import (
	"pg_sandbox/proto/api"
	credentialsservices "pg_sandbox/services/credentials_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GenerateSecretHandler(c *gin.Context) {

	var req api.RegenerateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	secret, err := credentialsservices.GenerateSecret(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to generate secret", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Secret generated!", gin.H{
		"secret": secret,
	})

}

func GenerateOAuthSignatureHandler(c *gin.Context) {

	var req api.GenerateOAuthSignatureRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	signature, err := credentialsservices.GenerateOAuthSignature(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to generate signature", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Signature generated!", gin.H{
		"signature": signature,
	})

}

func GetAPICredentialsHandler(c *gin.Context) {

	userID := c.Param("id")

	credentials, err := credentialsservices.GetAPICredentials(userID)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to retrieve credentials", err.Error())
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Credentials"), gin.H{
		"credentials": credentials,
	})

}

func SetPinHandler(c *gin.Context) {

	var req api.SetPinRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, 400, utils.FailBind)
		return
	}

	err := credentialsservices.SetPin(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to set/update pin", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Pin Updated successfully")

}
