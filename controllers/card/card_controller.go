package card

import (
	"log/slog"
	"pg_sandbox/proto/card"
	cardservices "pg_sandbox/services/card_services"
	commonservices "pg_sandbox/services/common_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func MakeCardRequestHandler(c *gin.Context) {

	xClientId := c.GetHeader("X-Client-ID")
	xTransactionRef := c.GetHeader("X-Transaction-Ref")

	commonservices.CheckEssentialHeaders(c)

	xCallbackUrl := c.GetHeader("X-CALLBACK-URL")

	if c.IsAborted() {
		return
	}

	var req card.CardRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		utils.Log(slog.LevelError, "Error", "Invalid Request Body", "Client ID", xClientId)

		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid Request Body.",
			"errors": gin.H{
				"Incomplete Body": []string{"The JSON body is missing"},
			},
		})

		return
	}

	cardservices.InitiateCardPayment(c, xClientId, xTransactionRef, xCallbackUrl, &req)
}

func Send3DsCodeHandler(c *gin.Context) {

	var req card.RequestCode

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	err = cardservices.SendCodeAccountHolder(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Failed to send OTP code, try again later", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "The OTP code has been sent please check your email")

}

func Verify3DsCodeHandler(c *gin.Context) {

	var req card.VerifyCode

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	err = cardservices.VerifyCardCode(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Failed to verify OTP code, try again later", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "The OTP code has been verified")

}
