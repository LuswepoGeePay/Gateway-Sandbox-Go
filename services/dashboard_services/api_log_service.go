package dashboardservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
	"time"
)

func GetApiStatistics() (*dashboard.APIStatisticsResponse, error) {

	var totalRequests int64

	var totalRequestsToday int64

	var errorRate int64
	now := time.Now()
	//users
	requestsQuery := config.DB.Model(&models.APILogs{})

	err := requestsQuery.Count(&totalRequests).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count requests")
	}

	//active
	errorQuery := config.DB.Model(&models.APILogs{})
	errorQuery = errorQuery.Where("status = ?", "failed")

	err = errorQuery.Count(&errorRate).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count errors")
	}

	location, _ := time.LoadLocation("Africa/Lusaka") // adjust based on deployment

	// Start of Today at 01:00 AM
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, location)
	endOfToday := startOfToday.AddDate(0, 0, 1)

	// Query
	newTodayQuery := config.DB.Model(&models.APILogs{}).
		Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday)

	err = newTodayQuery.Count(&totalRequestsToday).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count requests for today")
	}
	return &dashboard.APIStatisticsResponse{
		RequestsToday: int32(totalRequestsToday),
		Requests:      int32(totalRequests),
		ErrorRate:     int32(errorRate),
	}, nil
}

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
