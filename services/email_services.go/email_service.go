package services

import (
	"bytes"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/mail"
	"pg_sandbox/utils"

	"html/template"

	"gopkg.in/gomail.v2"
)

func SendCodeMail(req *mail.SendMailRequest) error {

	emailTemplate, err := template.ParseFiles("./email_template.html")

	if err != nil {
		return err
	}

	code := utils.GenerateSixDigitCode()

	var body bytes.Buffer

	err = emailTemplate.Execute(&body, struct{ Name string }{Name: code})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "luswepo17@gmail.com")
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "luswepo17@gmail.com", "dqiemeknokcbuexh")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	tx := config.DB.Begin()

	var user models.User

	err = config.DB.Where("email = ?", req.To).First(&user).Error

	if err != nil {
		return utils.CapitalizeError("failed to find user with that email")
	}

	updates := map[string]interface{}{}

	updates["otp_code"] = code
	updates["status"] = "active"

	if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// func SendGoMailWithAttachment(req *mail.SendMailWithAttachmentRequest, attachmentName string) error {
// 	// Create the "uploads" folder if it doesn't exist

// 	baseDir := "Files/Emails"

// 	var filePath string
// 	var err error

// 	// Create the email message
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", req.From)
// 	m.SetHeader("To", req.To)
// 	m.SetHeader("Subject", req.Subject)
// 	m.SetBody("text/plain", req.Body)

// 	if len(req.Attachment) > 0 {
// 		filePath, err = utils.SavePhoto(baseDir, uuid.New().String(), attachmentName, req.Attachment)
// 		if err != nil {
// 			return err
// 		}
// 		m.Attach(filePath)
// 	}

// 	// Attach the file
// 	m.Attach(filePath)

// 	// Set up the SMTP dialer
// 	d := gomail.NewDialer("smtp.gmail.com", 587, "luswepo17@gmail.com", "dqiemeknokcbuexh")
// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}

// 	return nil

// }

func VerifyEmail(req *mail.VerifyEmail) error {

	var user models.User

	tx := config.DB.Begin()

	err := tx.Where("email = ?", req.Email).First(&user).Error

	if err != nil {
		return utils.CapitalizeError("failed to find user")
	}

	updates := map[string]interface{}{}

	if req.Code == user.OtpCode {
		updates["email_verified"] = true
	}

	if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func SendChangePasswordCode(req *mail.SendMailRequest) error {

	// resetTemp := os.Getenv("RESET_PASSWORD_TEMPLATE_FILE")
	// if resetTemp == "" {
	// 	return utils.CapitalizeError("RESET_PASSWORD_FILE environment variable is not set")
	// }

	emailTemplate, err := template.ParseFiles("../reset_password.html")

	if err != nil {
		return err
	}

	var body bytes.Buffer

	code := utils.GenerateSixDigitCode()

	err = emailTemplate.Execute(&body, struct{ Name string }{Name: code})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "luswepo17@gmail.com")
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", "Password Reset!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "luswepo17@gmail.com", "dqiemeknokcbuexh")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	tx := config.DB.Begin()

	updates := map[string]interface{}{}

	updates["otp_code"] = code

	if err := tx.Model(&models.User{}).Where("id = ?", req.UserID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
