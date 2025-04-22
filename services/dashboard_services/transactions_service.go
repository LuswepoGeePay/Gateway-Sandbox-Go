package dashboardservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/dashboard"
	"pg_sandbox/utils"
	"time"
)

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
