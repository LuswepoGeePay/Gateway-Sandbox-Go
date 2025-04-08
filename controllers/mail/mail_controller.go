package mail

import (
	"pg_sandbox/proto/mail"
	services "pg_sandbox/services/email_services.go"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func SendMailHandler(c *gin.Context) {

	var req mail.SendMailRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, utils.FailBind, err.Error())
		return
	}

	err = services.SendCodeMail(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Failed to send email, try again later", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Sent Email, check your inbox!")

}

// func SendMailWithAttachmentHandler(c *gin.Context) {
// 	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
// 		utils.RespondWithError(c, 400, utils.UnParse, err.Error())
// 		return
// 	}

// 	emailData := c.Request.FormValue("email")

// 	if emailData == "" {
// 		utils.RespondWithError(c, 400, utils.MissData, "email data is required!")
// 		return
// 	}

// 	var req mail.SendMailWithAttachmentRequest
// 	if err := protojson.Unmarshal([]byte(emailData), &req); err != nil {
// 		utils.RespondWithError(c, 400, "invalid email data", err.Error())
// 		return
// 	}

// 	attachmentBytes, attachmentName := utils.ProcessPhoto(c, "attachment")

// 	req.Attachment = attachmentBytes

// 	err := services.SendGoMailWithAttachment(&req, attachmentName)

// 	if err != nil {
// 		utils.RespondWithError(c, 400, "Failed to send email, try again later", err.Error())
// 		return
// 	}

// 	// Respond with success
// 	utils.RespondWithSuccess(c, "Sent Email with Attachment, check your inbox!")
// }

func VerifyCode(c *gin.Context) {

	var req mail.VerifyEmail

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(c, 400, "Unable to verify email", err.Error())
		return
	}

	err = services.VerifyEmail(&req)

	if err != nil {
		utils.RespondWithError(c, 400, "Unable to verify email", err.Error())
		return
	}

	utils.RespondWithSuccess(c, "Email verified successfully!")
}
