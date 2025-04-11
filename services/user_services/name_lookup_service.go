package userservices

import (
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func NameLookUp(c *gin.Context, phoneNumber string) {

	network, err := utils.GetNetworkProvider(phoneNumber)

	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"phone_number": []string{"Invalid Phone number"},
			},
		})
		return

	}

	if network == "mtn" {
		c.JSON(200, gin.H{
			"code":    "200",
			"status":  "success",
			"message": "Name lookup completed successfully.",
			"data": gin.H{
				"status":       "success",
				"provider":     "MTN",
				"phone_number": phoneNumber,
				"names":        "John MTN Doe",
			},
		})
		return
	}

	if network == "airtel" {
		c.JSON(200, gin.H{
			"code":    "200",
			"status":  "success",
			"message": "Name lookup completed successfully.",
			"data": gin.H{
				"status":       "success",
				"provider":     "Airtel",
				"phone_number": phoneNumber,
				"names":        "Alice Airtel Bob",
			},
		})
		return
	}
	if network == "zamtel" {
		c.JSON(200, gin.H{
			"code":    "200",
			"status":  "success",
			"message": "Name lookup completed successfully.",
			"data": gin.H{
				"status":       "success",
				"provider":     "Zamtel",
				"phone_number": phoneNumber,
				"names":        "Nagato Zamtel Gato",
			},
		})
		return
	}

}
