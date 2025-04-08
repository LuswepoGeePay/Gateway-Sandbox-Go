package userservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/user"
	"pg_sandbox/utils"

	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(req *user.ResetPassword) error {

	if req.Password != req.ConfirmPassword {
		return utils.CapitalizeError("password do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return utils.CapitalizeError("unable to hash password.")
	}

	updates := map[string]interface{}{}

	var user models.User

	if req.LoggedIn {
		err = config.DB.Where("id =  ?", req.UserId).First(&user).Error

		if err != nil {
			return utils.CapitalizeError("unable to find user")
		}
	}

	if !req.LoggedIn {
		err = config.DB.Where("email =  ?", req.Email).First(&user).Error

		if err != nil {
			return utils.CapitalizeError("unable to find user")
		}
	}

	if req.Password != "" {
		updates["password"] = hashedPassword
	}

	tx := config.DB.Begin()

	if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
