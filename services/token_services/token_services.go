package tokenservices

import (
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	pbToken "pg_sandbox/proto/token"
	"pg_sandbox/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("SuperSecretKeyForARobustSystem")

func GenerateOAuthToken(req *pbToken.TokenRequest) (*pbToken.TokenResponse, error) {

	if req.GrantType != "client_credentials" {
		utils.Log(slog.LevelError, "wrong grant type", "data", gin.H{
			"grant_type": req.GrantType,
		})
		return nil, utils.CapitalizeError("invalid Grant Type")
	}

	var apiKey models.ApiKeys

	tx := config.DB.Begin()
	if err := tx.Where("client_id = ? AND client_secret = ? AND is_active = ?", req.ClientId, req.ClientSecret, true).First(&apiKey).Error; err != nil {
		utils.Log(slog.LevelError, "unable to find user", "data", gin.H{
			"client_secret": req.ClientSecret,
			"client_id":     req.ClientId,
		})
		return nil, utils.CapitalizeError("user not found")
	}

	expiresAt := time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   apiKey.UserID,
		"client_id": req.ClientId,
		"exp":       expiresAt,
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		utils.Log(slog.LevelError, "unable to generate authorization token", "data", gin.H{
			"user_id":   apiKey.UserID,
			"client_id": req.ClientId,
		})
		return nil, utils.CapitalizeError("failed to generate token")
	}

	return &pbToken.TokenResponse{
		TokenType:   "Bearer",
		ExpiresIn:   int32(expiresAt),
		AccessToken: tokenString,
	}, nil
}

func ValidateOAuthToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.CapitalizeError("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return utils.CapitalizeError("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return utils.CapitalizeError("invalid token claims")
	}

	return nil
}
