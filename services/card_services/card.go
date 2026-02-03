package cardservices

import (
	"bytes"
	"context"
	"crypto/tls"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/card"
	"pg_sandbox/services/logs"
	"pg_sandbox/utils"
	"strings"
	"text/template"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func InitiateCardPayment(c *gin.Context, xClientId string, xTransactionRef string, xCallbackUrl string, req *card.CardRequest) {

	start := time.Now()

	var existingClientID models.ApiKeys

	result := config.DB.Where("client_id = ?", xClientId).First(&existingClientID)

	if result.Error != nil {

		utils.Log(slog.LevelError, "❌Error", "unable to initiate card payment, client ID is invalid", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})
		c.JSON(402, gin.H{
			"code":    402,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Client-ID": []string{"The selected x-client-id is invalid."},
			},
		})
		return
	}
	if req.CardNumber == "" || req.Amount == "" || xTransactionRef == "" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/card/payment", "POST", "failed", strconv.FormatInt(elapsed, 10))
		utils.Log(slog.LevelError, "❌Error", "unable to initiate card payment, card details or amount is empty", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})
		c.JSON(402, gin.H{
			"code":    402,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Card-Number":     []string{"Card number cannot be empty."},
				"X-Amount":          []string{"Amount cannot be empty."},
				"X-Transaction-Ref": []string{"Transaction reference cannot be empty."},
			},
		})
		return
	}

	if strings.TrimSpace(xTransactionRef) == "" {
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/card/payment", "POST", "failed", strconv.FormatInt(elapsed, 10))

		utils.Log(slog.LevelError, "❌Error", "unable to initate request to pay, transaction reference is empty", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})

		c.JSON(402, gin.H{
			"code":    402,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Transaction-Ref": []string{"Transaction reference cannot be empty."},
			},
		})
		return
	}

	var existingTransaction models.Transactions

	result = config.DB.Where("reference = ?", xTransactionRef).First(&existingTransaction)

	if result.Error == nil { // Only return an error if the transaction exists
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, existingClientID.UserID.String(), "/v1/card/payment", "POST", "failed", strconv.FormatInt(elapsed, 10))

		utils.Log(slog.LevelError, "❌Error", "unable to initiate card payment, transaction reference has already been taken", "data", gin.H{
			"client_id":       xClientId,
			"transaction_ref": xTransactionRef,
			"request_body":    req,
		})

		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"X-Transaction-Ref": []string{"The x-transaction-ref has already been taken."},
			},
		})
		return
	}

	transaction := models.Transactions{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(req.UserId),
		Reference: xTransactionRef,
		Amount:    req.Amount,
		Status:    "pending",
		Customer:  req.CardHolderName,
		Channel:   "card",
		Type:      "collection",
		Date:      time.Now(),
	}

	tx := config.DB.Begin()
	result = tx.Create(&transaction)

	if result.Error != nil {
		tx.Rollback()
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, req.UserId, "/v1/card/payment", "POST", "failed", strconv.FormatInt(elapsed, 10))
		utils.Log(slog.LevelError, "❌Error", "unable to initiate card payment, unable to create transaction", "req", req)
		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "System Error",
			"errors": gin.H{
				"Tranaction Error": []string{"Failed to create transaction"},
			},
		})
		return
	}

	urlId := uuid.New()

	cardUrl := req.BaseUrl + urlId.String()

	cardUrlModel := models.CardUrls{
		ID:            uuid.New(),
		Url:           cardUrl,
		UserID:        uuid.MustParse(req.UserId),
		Code:          "",
		Verified:      false,
		TransactionId: xTransactionRef,
	}

	result = tx.Create(&cardUrlModel)

	if result.Error != nil {
		tx.Rollback()
		elapsed := time.Since(start).Milliseconds()
		logs.LogApiCall(c, req.UserId, "/v1/card/payment", "POST", "failed", strconv.FormatInt(elapsed, 10))
		utils.Log(slog.LevelError, "❌Error", "unable to save card url", "req", req)
		c.JSON(422, gin.H{
			"code":    422,
			"status":  "error",
			"message": "System Error",
			"errors": gin.H{
				"Tranaction Error": []string{"Failed to save card url"},
			},
		})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{
		"status":                "success",
		"message":               "Card transaction created",
		"url":                   cardUrl,
		"transaction_reference": xTransactionRef,
	})

}

func SendCodeAccountHolder(req *card.RequestCode) error {

	var url models.CardUrls

	if err := config.DB.Where("user_id = ?  AND transaction_id = ?", req.UserId, req.TransactionReference).Find(&url).Error; err != nil {
		return err
	}

	emailTemplate, err := template.ParseFiles("./send_code.html")

	if err != nil {
		return err
	}

	code := utils.GenerateSixDigitCode()

	var body bytes.Buffer

	err = emailTemplate.Execute(&body, struct{ Name string }{Name: code})

	if err != nil {
		return err
	}

	var user models.User

	if err := config.DB.Where("id = ?", req.UserId).Find(&user).Error; err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "gpgsnoreply@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "OTP Code")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		"smtp.gmail.com",        // MAIL_HOST
		465,                     // MAIL_PORT
		"gpgsnoreply@gmail.com", // MAIL_USERNAME
		"pbeb pnvy kdkt ykze",   // MAIL_PASSWORD
	)
	d.SSL = true
	d.TLSConfig = &tls.Config{ServerName: "smtp.gmail.com"}

	// Add timeout
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	tx := config.DB.Begin()

	// var user models.User

	// err = config.DB.Where("email = ?", req.To).First(&user).Error

	// if err != nil {
	// 	return utils.CapitalizeError("failed to find user with that email")
	// }

	updates := map[string]interface{}{}

	updates["code"] = code

	if err := tx.Model(&models.CardUrls{}).Where("id = ?", url.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func VerifyCardCode(req *card.VerifyCode) error {

	var url models.CardUrls

	tx := config.DB.Begin()

	err := tx.Where("transaction_id = ? AND  user_id = ?", req.TransactionReference, req.UserId).First(&url).Error

	if err != nil {
		return utils.CapitalizeError("failed to find transaction")
	}

	updates := map[string]interface{}{}

	if req.Code != url.Code {
		return utils.CapitalizeError("Invalid code, please try again")
	}

	if req.Code == url.Code {
		updates["verified"] = true
	}

	if err := tx.Model(&models.CardUrls{}).Where("id = ?", url.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	var transaction models.Transactions
	if err := config.DB.Where("reference = ?", req.TransactionReference).First(&transaction).Error; err != nil {
		return err
	}

	transaction.Status = "successful" // Update status to completed after OTP verification
	if err := config.DB.Save(&transaction).Error; err != nil {
		return err
	}

	tx.Commit()

	return nil
}
