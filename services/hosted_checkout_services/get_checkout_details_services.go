package hostedcheckoutservices

import (
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/hcheckout"
	"pg_sandbox/utils"
)

func GetCheckoutDetails(checkoutID string) (*hcheckout.HCheckout, error) {

	var checkout models.CheckOutUrls

	result := config.DB.Where("id = ?", checkoutID).Find(&checkout)

	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find url")
	}

	return &hcheckout.HCheckout{
		CheckoutUrl: checkout.GeneratedUrl,
		OrderId:     checkout.OrderID,
		Customer: &hcheckout.Customer{
			Name:  checkout.CustomerName,
			Email: checkout.CustomerEmail,
		},
		ReturnUrl: checkout.ReturnUrl,
		Amount:    checkout.Amount,
	}, nil
}
