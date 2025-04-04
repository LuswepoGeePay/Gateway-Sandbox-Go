package auth

import (
	"net/http"
	"pg_sandbox/proto/token"
	tokenservices "pg_sandbox/services/token_services"

	"github.com/gin-gonic/gin"
)

func AuthorizationHandler(c *gin.Context) {

	var req token.TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := tokenservices.GenerateOAuthToken(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	c.JSON(200, response)
}
