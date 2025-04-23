package dashboardservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
)

func GetActivity(req *dashboard.GetActivityrequests) (*dashboard.GetActivityResponse, error) {

	var activities []models.ActivityLogs

	tx := config.DB.Begin()
	query := tx.Model(&models.ActivityLogs{})

	var totalActivities int64
	err := query.Count(&totalActivities).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count activities")
	}

	totalPages := int32((totalActivities + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&activities).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve activities")
	}

	pbActivities := make([]*dashboard.Activity, len(activities))

	for i, activity := range activities {

		pbActivities[i] = &dashboard.Activity{
			Id:        activity.ID.String(),
			UserId:    activity.UserID.String(),
			Entity:    activity.Entity,
			Action:    activity.Action,
			EntityId:  activity.EntityID,
			Ipaddress: activity.IPAddress,
		}
	}

	return &dashboard.GetActivityResponse{
		Activity:    pbActivities,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}
