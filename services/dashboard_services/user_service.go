package dashboardservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
)

func GetUsers(req *dashboard.GetUsersRequest) (*dashboard.GetUsersResponse, error) {

	var users []models.User

	tx := config.DB.Begin()
	query := tx.Model(&models.User{})

	var totalUsers int64
	err := query.Count(&totalUsers).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count users")
	}

	totalPages := int32((totalUsers + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&users).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve users")
	}

	pbUsers := make([]*dashboard.User, len(users))

	for i, user := range users {

		pbUsers[i] = &dashboard.User{
			Id:       user.ID.String(),
			Fullname: user.Fullname,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role.Name,
			Status:   user.Status,
		}
	}

	return &dashboard.GetUsersResponse{
		User:        pbUsers,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}

func GetUserStatistics() (*dashboard.UserStatisticsResponse, error) {
	var totalMerchants int64
	var totalUsers int64
	var totalActiveUsers int64
	var totalInActiveUsers int64
	var totalAdmins int64
	//users
	usersQuery := config.DB.Model(&models.User{})

	err := usersQuery.Count(&totalUsers).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//active
	activeQuery := config.DB.Model(&models.User{})
	activeQuery = activeQuery.Where("status = ?", "active")

	err = activeQuery.Count(&totalActiveUsers).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//inactive
	inactiveQuery := config.DB.Model(&models.User{})
	inactiveQuery = inactiveQuery.Where("status = ?", "inactive")

	err = inactiveQuery.Count(&totalInActiveUsers).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//MERCHANTs

	var role models.Role
	result := config.DB.Where("name = ?", "merchant").First(&role)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role merchant.")
	}

	merchantQuery := config.DB.Model(&models.User{})
	merchantQuery = merchantQuery.Where("role_id = ?", role.ID)

	err = merchantQuery.Count(&totalMerchants).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//Admin

	var adminRole models.Role
	result = config.DB.Where("name = ?", "admin").First(&adminRole)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role admin.")
	}

	adminQuery := config.DB.Model(&models.User{})
	adminQuery = adminQuery.Where("role_id = ?", adminRole.ID)

	err = adminQuery.Count(&totalAdmins).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	return &dashboard.UserStatisticsResponse{
		AdminUsers:    int32(totalAdmins),
		ActiveUsers:   int32(totalActiveUsers),
		InactiveUsers: int32(totalInActiveUsers),
		MerchantUsers: int32(totalMerchants),
		Users:         int32(totalUsers),
	}, nil

}
