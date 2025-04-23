package dashboardservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
	"time"

	"github.com/google/uuid"
)

func GetMerchants(req *dashboard.GetUsersRequest) (*dashboard.GetUsersResponse, error) {

	var users []models.User

	tx := config.DB.Begin()
	query := tx.Model(&models.User{})

	var role models.Role
	result := tx.Where("name = ?", "merchant").First(&role)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role.")
	}

	var totalUsers int64
	err := query.Where("role_id = ?", role.ID).Count(&totalUsers).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count merchants")
	}

	totalPages := int32((totalUsers + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&users).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve merchants")
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

func GetMerchantStatistics() (*dashboard.MerchantStatisticsResponse, error) {
	var totalMerchants int64
	var totalActiveMerchants int64
	var totalInActiveMerchants int64
	var newToday int64
	var newThisWeek int64
	var newThisMonth int64

	var role models.Role
	result := config.DB.Where("name = ?", "merchant").First(&role)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role.")
	}

	//active
	activeQuery := config.DB.Model(&models.User{})
	activeQuery = activeQuery.Where("status = ?  AND role_id = ? ", "active", role.ID)

	err := activeQuery.Count(&totalActiveMerchants).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//inactive
	inactiveQuery := config.DB.Model(&models.User{})
	inactiveQuery = inactiveQuery.Where("status = ? AND role_id = ? ", "inactive", role.ID)

	err = inactiveQuery.Count(&totalInActiveMerchants).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	//MERCHANTs

	merchantQuery := config.DB.Model(&models.User{})
	merchantQuery = merchantQuery.Where("role_id = ?", role.ID)

	err = merchantQuery.Count(&totalMerchants).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users")
	}

	now := time.Now()

	// This Week
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // make Sunday = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	// This Month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// New This Week
	weekQuery := config.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", startOfWeek, endOfWeek).
		Where("role_id = ?", role.ID)

	err = weekQuery.Count(&newThisWeek).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users this week")
	}

	// New This Month
	monthQuery := config.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", startOfMonth, endOfMonth).
		Where("role_id = ?", role.ID)

	err = monthQuery.Count(&newThisMonth).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users this month")
	}

	//New Today

	location, _ := time.LoadLocation("Africa/Lusaka") // adjust based on deployment

	// Start of Today at 01:00 AM
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, location)
	endOfToday := startOfToday.AddDate(0, 0, 1)

	// Query
	newTodayQuery := config.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday).
		Where("role_id = ?", role.ID)

	err = newTodayQuery.Count(&newToday).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count users for today")
	}

	return &dashboard.MerchantStatisticsResponse{
		ActiveMerchants:   int32(totalActiveMerchants),
		Merchants:         int32(totalMerchants),
		InactiveMerchants: int32(totalInActiveMerchants),
		NewToday:          int32(newToday),
		NewMonth:          int32(newThisMonth),
		NewWeek:           int32(newThisWeek),
	}, nil

}

func GetTopMerchants() (*dashboard.TopMerchantsResponse, error) {

	type MerchantInfo struct {
		MerchantID uuid.UUID
		Fullname   string
		Count      int
	}

	var merchantRole models.Role
	err := config.DB.Where("name = ?", "merchant").First(&merchantRole).Error
	if err != nil {
		// handle error, e.g., role not found
		return nil, utils.CapitalizeError("unable to fetch roles")
	}

	var topMerchants []MerchantInfo

	if err := config.DB.Table("transactions").
		Select("users.id as merchant_id, users.fullname, COUNT(transactions.id) as count").
		Joins("LEFT JOIN users ON users.id = transactions.user_id").
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.role_id = ?", merchantRole.ID). // Replace with actual merchant role UUID or dynamic role ID fetch
		Group("users.id, users.fullname").
		Order("count DESC").
		Limit(3).
		Scan(&topMerchants).Error; err != nil {
		return nil, utils.CapitalizeError("unable to fetch top merchants")
	}

	resp := &dashboard.TopMerchantsResponse{}
	if len(topMerchants) > 0 {
		resp.MerchantOne = topMerchants[0].Fullname
		resp.MerchantOneCount = int32(topMerchants[0].Count)
	}
	if len(topMerchants) > 1 {
		resp.MerchantTwo = topMerchants[1].Fullname
		resp.MerchantTwoCount = int32(topMerchants[1].Count)
	}
	if len(topMerchants) > 2 {
		resp.MerchantThree = topMerchants[2].Fullname
		resp.MerchantThreeCount = int32(topMerchants[2].Count)
	}

	return &dashboard.TopMerchantsResponse{
		MerchantOne:        resp.MerchantOne,
		MerchantTwo:        resp.MerchantTwo,
		MerchantThree:      resp.MerchantThree,
		MerchantOneCount:   resp.MerchantOneCount,
		MerchantTwoCount:   resp.MerchantTwoCount,
		MerchantThreeCount: resp.MerchantThreeCount,
	}, nil
}
