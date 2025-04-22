package dashboardservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
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
