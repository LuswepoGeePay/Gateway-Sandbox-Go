package authservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	"pg_sandbox/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req *auth.RegisterRequest) (*string, error) {

	tx := config.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, utils.CapitalizeError("unable to hash password")
	}

	var role models.Role
	result := tx.Where("name = ?", req.Role).First(&role)
	if result.Error != nil {
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
		tx.Rollback()
		return nil, utils.CapitalizeError(result.Error.Error())
	}

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
		return nil, utils.CapitalizeError("unable to create API keys")
	}

	tx.Commit()

	userIdStr := userId.String()

	return &userIdStr, nil
}
