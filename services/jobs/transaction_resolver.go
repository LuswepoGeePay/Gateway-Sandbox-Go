package jobs

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	disbursementservices "pg_sandbox/services/disbursement_services"
	"pg_sandbox/utils"
	"time"
)

func StartTransactionResolver() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	go func() {
		for range ticker.C {
			ResolvePendingTransactions()
		}
	}()
}

func ResolvePendingTransactions() {
	var transactions []models.Transactions

	// Find all pending collection transactions
	result := config.DB.Where("status = ? AND type = ?", "pending", "collection").Find(&transactions)
	if result.Error != nil {
		utils.Log(slog.LevelError, "Error fetching pending transactions", "error", result.Error)
		return
	}

	for _, tx := range transactions {
		// Update status to successful
		tx.Status = "successful"
		if err := config.DB.Save(&tx).Error; err != nil {
			utils.Log(slog.LevelError, "Error updating transaction status to successful", "tx_id", tx.ID, "error", err)
			continue
		}

		utils.Log(slog.LevelInfo, "Transaction auto-resolved to successful", "tx_id", tx.ID, "ref", tx.Reference)

		// Fire callback if callback URL is provided
		if tx.CallbackUrl != "" {
			tCode := utils.GenerateTenDigitCode()
			go disbursementservices.CallbackHandler(tx.CallbackUrl, models.CallbackPayload{
				Code:    200,
				Status:  "successful",
				Message: "Transaction has been successfully processed and settled.",
				Data: models.CallbackPayloadData{
					TransactionReference: tx.Reference,
					ExternalReference:    tCode,
					Customer:             tx.Customer,
					Amount:               tx.Amount,
				},
			})
		}
	}
}
