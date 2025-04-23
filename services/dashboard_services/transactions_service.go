package dashboardservices

import (
	"fmt"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
	"strconv"
	"strings"
	"time"
)

func GetTransactionStatistics() (*dashboard.TransactionStatisticsResponse, error) {

	var totalTransactions int64

	var totalSuccessful int64

	var totalFailed int64
	var totalPending int64
	var amountTotal float64

	//users
	transactionsQuery := config.DB.Model(&models.Transactions{})

	err := transactionsQuery.Count(&totalTransactions).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count transactions")

	}

	//successful
	successQuery := config.DB.Model(&models.Transactions{})
	successQuery = successQuery.Where("status = ? OR status = ?", "successful", "completed")

	err = successQuery.Count(&totalSuccessful).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count successful")
	}

	//pending
	pendingQuery := config.DB.Model(&models.Transactions{})
	pendingQuery = pendingQuery.Where("status = ?", "pending")

	err = pendingQuery.Count(&totalPending).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count pending transactions")
	}

	//Errored
	errorQuery := config.DB.Model(&models.Transactions{})
	errorQuery = errorQuery.Where("status = ?", "failed")

	err = errorQuery.Count(&totalFailed).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to count failed transactions")
	}

	//amount

	var transactions []models.Transactions
	amountTotal = 0

	amountQuery := config.DB.Model(&models.Transactions{})
	err = amountQuery.Select("amount").Find(&transactions).Error
	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to fetch transactions for amount sum")
	}

	for _, tx := range transactions {
		amt, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			utils.Log(slog.LevelWarn, "Amount Parse Error", fmt.Sprintf("txID: %v, amount: %v", tx.ID, tx.Amount))
			continue // skip faulty entries
		}
		amountTotal += amt
	}

	return &dashboard.TransactionStatisticsResponse{
		Transactions: int32(totalTransactions),
		Successful:   int32(totalSuccessful),
		Pending:      int32(totalPending),
		Failed:       int32(totalFailed),
		TotalAmount:  int32(amountTotal),
	}, nil
}

func GetTransactions(req *dashboard.GetTransactionsRequest) (*dashboard.GetTransactionsResponse, error) {
	var transactions []models.Transactions

	query := config.DB.Model(&models.Transactions{})

	var totalTransactions int64
	err := query.Count(&totalTransactions).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count transactions")
	}

	totalPages := int32((totalTransactions + int64(req.PageSize) - 1) / int64(req.PageSize))

	// if req.SearchQuery != "" {
	// 	searchValue := "%" + strings.ToLower(req.SearchQuery) + "%"
	// 	query = query.Where("LOWER(transaction_name) LIKE ?", searchValue)
	// }

	// if req.Status != "" {
	// 	query = query.Where("LOWER(status) LIKE ?", req.Status)
	// }

	//calcualate offset
	offset := (req.Page - 1) * req.PageSize

	err = query.Order("transactions.created_at DESC").Limit(int(req.PageSize)).
		Offset(int(offset)).Find(&transactions).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve evnts")
	}

	pbtransactions := make([]*dashboard.Transaction, len(transactions))

	for i, transaction := range transactions {
		pbtransactions[i] = &dashboard.Transaction{
			Id:        transaction.ID.String(),
			Reference: transaction.Reference,
			Amount:    transaction.Amount,
			Status:    transaction.Status,
			Customer:  transaction.Customer,
			Channel:   transaction.Channel,
			Type:      transaction.Type,
			Narration: transaction.Narration,
			Date:      transaction.Date.Format(time.RFC3339),
		}
	}
	return &dashboard.GetTransactionsResponse{
		Transaction: pbtransactions,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}

func GetTransactionChannelStats() (*dashboard.TransactionChannelsResponse, error) {
	type ChannelStat struct {
		Channel string
		Count   int64
	}

	var results []ChannelStat
	err := config.DB.
		Model(&models.Transactions{}).
		Select("channel, COUNT(*) as count").
		Group("channel").
		Scan(&results).Error

	if err != nil {
		utils.Log(slog.LevelError, "Error", err.Error())
		return nil, utils.CapitalizeError("failed to fetch transaction channel stats")
	}

	resp := &dashboard.TransactionChannelsResponse{}
	for _, stat := range results {
		switch strings.ToLower(stat.Channel) {
		case "mtn":
			resp.Mtn = int32(stat.Count)
		case "zamtel":
			resp.Zamtel = int32(stat.Count)
		case "airtel":
			resp.Airtel = int32(stat.Count)
		}
	}
	return resp, nil
}
