package adminservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	"pg_sandbox/utils"
)

func GetMerchants(req *auth.GetUsersRequest) (*auth.GetUsersResponse, error) {

	var authModel []models.User
	offset := (req.Page - 1) * req.PageSize

	tx := config.DB.Begin()

	var totalUsers int64

	query := tx.Model(&models.User{}).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ? ", "merchant")

	// Count total users matching the criteria
	if err := query.Count(&totalUsers).Error; err != nil {
		utils.Log(slog.LevelError, "Failed to count users", "error", err.Error())
		tx.Rollback()
		return nil, utils.CapitalizeError(err.Error())
	}

	totalPages := int32((totalUsers + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Fetch users with pagination
	if err := query.Limit(int(req.PageSize)).Offset(int(offset)).Find(&authModel).Error; err != nil {
		utils.Log(slog.LevelError, "Failed to retrieve users", "error", err.Error())
		tx.Rollback()
		return nil, utils.CapitalizeError(err.Error())
	}

	tx.Commit()

	var pbUser []*auth.User
	for _, authI := range authModel {
		pbUser = append(pbUser, &auth.User{
			Id:       authI.ID.String(),
			Fullname: authI.Fullname,
			Email:    authI.Email,
			Phone:    authI.Phone,
			Role:     authI.Role.Name,
			Status:   authI.Status,
		})
	}

	return &auth.GetUsersResponse{
		User:        pbUser,
		TotalPages:  totalPages,
		HasMore:     req.Page < totalPages,
		CurrentPage: req.Page,
	}, nil
}
