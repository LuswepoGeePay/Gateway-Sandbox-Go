package dashboardservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
)

func GetOverviewCardsInfo() (*dashboard.OverviewCardInfoResponse, error) {

	var totalMerchants int64
	var totalUsers int64
	var totalAPIReqs int64

	//users
	usersQuery := config.DB.Model(&models.User{})

	err := usersQuery.Count(&totalUsers).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//api requests
	apiQuery := config.DB.Model(&models.APILogs{})

	err = apiQuery.Count(&totalAPIReqs).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//MERCHANTs

	var role models.Role
	result := config.DB.Where("name = ?", "merchant").First(&role)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role.")
	}

	merchantQuery := config.DB.Model(&models.User{})
	merchantQuery = merchantQuery.Where("role_id = ?", role.ID)

	err = merchantQuery.Count(&totalMerchants).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	return &dashboard.OverviewCardInfoResponse{
		ApiRequests: int32(totalAPIReqs),
		Merchants:   int32(totalMerchants),
		Users:       int32(totalUsers),
	}, nil
}
