package userservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	pbUser "pg_sandbox/proto/user"
	"pg_sandbox/utils"

	"github.com/google/uuid"
)

func GetUserProfile(userId string) (*pbUser.User, error) {

	userID, err := uuid.Parse(userId)

	if err != nil {
		return nil, utils.CapitalizeError("invalid user UUID")
	}

	var user models.User

	result := config.DB.Where("id = ?", userID).Find(&user)

	if result.Error != nil {
		return nil, utils.CapitalizeError(result.Error.Error())
	}

	return &pbUser.User{
		Fullname:      user.Fullname,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Phone:         user.Phone,
		Status:        user.Status,
	}, nil
}
