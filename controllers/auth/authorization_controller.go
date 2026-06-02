package auth

import (
	"net/http"
	"pg_sandbox/proto/token"
	tokenservices "pg_sandbox/services/token_services"

	"github.com/gin-gonic/gin"
)

type TokenHTTPRequest struct {
	GrantType    string `form:"grant_type" json:"grant_type"`
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

func AuthorizationHandler(c *gin.Context) {

	var httpReq TokenHTTPRequest

	if err := c.ShouldBind(&httpReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	req := &token.TokenRequest{
		GrantType:    httpReq.GrantType,
		ClientId:     httpReq.ClientId,
		ClientSecret: httpReq.ClientSecret,
	}
	response, err := tokenservices.GenerateOAuthToken(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}
