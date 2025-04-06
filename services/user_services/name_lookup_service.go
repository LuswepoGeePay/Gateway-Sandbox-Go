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
			"message": "Names retrieved successfully",
			"data": gin.H{
				"names": "John MTN Doe",
			},
		})
		return
	}

	if network == "airtel" {
		c.JSON(200, gin.H{
			"code":    "200",
			"status":  "success",
			"message": "Names retrieved successfully",
			"data": gin.H{
				"names": "Alice Airtel Bob",
			},
		})
		return
	}
	if network == "zamtel" {
		c.JSON(200, gin.H{
			"code":    "200",
			"status":  "success",
			"message": "Names retrieved successfully",
			"data": gin.H{
				"names": "Nagato Zamtel Gato",
			},
		})
		return
	}

}
