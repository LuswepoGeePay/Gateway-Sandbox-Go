package dashboardservices

import (
	"fmt"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
	"time"
)

func GetApiStatistics() (*dashboard.APIStatisticsResponse, error) {

	type EndpointStat struct {
		Endpoint string
		Count    int64
	}

	var endpointStats []EndpointStat
	err := config.DB.
		Model(&models.APILogs{}).
		Select("endpoint, COUNT(*) as count").
		Group("endpoint").
		Order("count DESC").
		Limit(3).
		Scan(&endpointStats).Error

	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to fetch endpoint stats")
	}
	var totalRequests int64

	var totalRequestsToday int64

	var errorRate int64
	now := time.Now()
	//users
	requestsQuery := config.DB.Model(&models.APILogs{})

	err = requestsQuery.Count(&totalRequests).Error
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

	resp := &dashboard.APIStatisticsResponse{
		RequestsToday: int32(totalRequestsToday),
		Requests:      int32(totalRequests),
		ErrorRate:     int32(errorRate),
	}

	if len(endpointStats) > 0 {
		resp.Endpoint1 = endpointStats[0].Endpoint
		resp.Endpoint1Count = fmt.Sprintf("%d", endpointStats[0].Count)
	}
	if len(endpointStats) > 1 {
		resp.Endpoint2 = endpointStats[1].Endpoint
		resp.Endpoint2Count = fmt.Sprintf("%d", endpointStats[1].Count)
	}
	if len(endpointStats) > 2 {
		resp.Endpoint3 = endpointStats[2].Endpoint
		resp.Endpoint3Count = fmt.Sprintf("%d", endpointStats[2].Count)
	}
	return &dashboard.APIStatisticsResponse{
		RequestsToday:  resp.RequestsToday,
		Requests:       resp.Requests,
		ErrorRate:      resp.ErrorRate,
		Endpoint1:      resp.Endpoint1,
		Endpoint2:      resp.Endpoint2,
		Endpoint3:      resp.Endpoint3,
		Endpoint1Count: resp.Endpoint1Count,
		Endpoint2Count: resp.Endpoint2Count,
		Endpoint3Count: resp.Endpoint3Count,
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

	err = query.Preload("User").Order("api_logs.created_at DESC").Limit(int(req.PageSize)).
		Offset(int(offset)).Find(&apiLogs).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve evnts")
	}

	pbapiLogs := make([]*dashboard.API, len(apiLogs))

	for i, log := range apiLogs {
		pbapiLogs[i] = &dashboard.API{
			Id:           log.ID.String(),
			UserId:       log.User.Fullname,
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

func GetAPIResponseTimeStats() (*dashboard.APIResTimeStatisticsResponse, error) {
	type ResponseStat struct {
		Endpoint    string
		AverageTime float64
	}

	var stats []ResponseStat
	err := config.DB.
		Model(&models.APILogs{}).
		Select("endpoint, AVG(response_time) as average_time").
		Group("endpoint").
		Order("average_time DESC").
		Limit(3).
		Scan(&stats).Error

	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to fetch API response times")
	}

	resp := &dashboard.APIResTimeStatisticsResponse{}
	if len(stats) > 0 {
		resp.Endpoint1 = stats[0].Endpoint
		resp.Endpoint1Time = fmt.Sprintf("%.2f ms", stats[0].AverageTime)
	}
	if len(stats) > 1 {
		resp.Endpoint2 = stats[1].Endpoint
		resp.Endpoint2Time = fmt.Sprintf("%.2f ms", stats[1].AverageTime)
	}
	if len(stats) > 2 {
		resp.Endpoint3 = stats[2].Endpoint
		resp.Endpoint3Time = fmt.Sprintf("%.2f ms", stats[2].AverageTime)
	}
	return resp, nil
}
