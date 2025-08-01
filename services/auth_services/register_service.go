package authservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context, req *auth.RegisterRequest) (*string, error) {

	tx := config.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		utils.Log(slog.LevelError, "❌Error", "unable to generate hash from password", "data", gin.H{
			"email": req.Email,
			"error": err,
		})
		return nil, utils.CapitalizeError("unable to hash password")
	}

	var role models.Role
	result := tx.Where("name = ?", req.Role).First(&role)
	if result.Error != nil {
		utils.Log(slog.LevelError, "❌Error", "unable to find role", "data", gin.H{
			"email":    req.Email,
			"error":    err,
			"role":     req.Role,
			"phone":    req.Phone,
			"fullname": req.Fullname,
		})
		return nil, utils.CapitalizeError("unable to find role.")
	}

	userId := uuid.New()
	user := models.User{
		ID:       userId,
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
		Status:   "active",
	}

	result = tx.Create(&user)
	if result.Error != nil {
		utils.Log(slog.LevelError, "❌Error", "unable to create user", "data", gin.H{
			"email":    req.Email,
			"error":    err,
			"role":     req.Role,
			"phone":    req.Phone,
			"fullname": req.Fullname,
		})
		tx.Rollback()
		return nil, utils.CapitalizeError(result.Error.Error())
	}

	// logs.LogActivity(c, userId.String(), "Registration", "Auth", userId.String())

	clientID := uuid.New().String()
	clientSecret := uuid.New().String() // You can use a more secure approach to generate the secret

	apiKey := models.ApiKeys{
		ID:           uuid.New(),
		UserID:       user.ID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		IsActive:     true,
	}

	if err := tx.Create(&apiKey).Error; err != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌Error", "unable to create api keys", "data", gin.H{
			"email":         req.Email,
			"client_id":     clientID,
			"client_secret": clientSecret,
		})
		return nil, utils.CapitalizeError("unable to create API keys")
	}

	tx.Commit()

	userIdStr := userId.String()

	return &userIdStr, nil
}
