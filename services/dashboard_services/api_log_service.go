package dashboardservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
)

func GetAPIRequests(req *dashboard.GetAPIrequests) (*dashboard.GetAPIReqsResponse, error) {
	var apiLogs []models.APILogs

	query := config.DB.Model(&models.APILogs{})

	var totalApiLogs int64
	err := query.Count(&totalApiLogs).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count api Log entries")
	}

	totalPages := int32((totalApiLogs + int64(req.PageSize) - 1) / int64(req.PageSize))

	// if req.SearchQuery != "" {
	// 	searchValue := "%" + strings.ToLower(req.SearchQuery) + "%"
	// 	query = query.Where("LOWER(log_name) LIKE ?", searchValue)
	// }

	// if req.Status != "" {
	// 	query = query.Where("LOWER(status) LIKE ?", req.Status)
	// }

	//calcualate offset
	offset := (req.Page - 1) * req.PageSize

	err = query.Order("api_logs.created_at DESC").Limit(int(req.PageSize)).
		Offset(int(offset)).Find(&apiLogs).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve evnts")
	}

	pbapiLogs := make([]*dashboard.API, len(apiLogs))

	for i, log := range apiLogs {
		pbapiLogs[i] = &dashboard.API{
			Id:           log.ID.String(),
			UserId:       log.UserID.String(),
			Status:       log.Status,
			Endpoint:     log.Endpoint,
			Method:       log.Method,
			Ipaddress:    log.IPAddress,
			ResponseTime: log.ResponseTime,
		}
	}
	return &dashboard.GetAPIReqsResponse{
		Request:     pbapiLogs,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}
