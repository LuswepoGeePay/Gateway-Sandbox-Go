package authservices

import (
	"log/slog"
	"pg_sandbox/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("tyhsdndfuadbajsddoewkmdiedwnnpewesedrftgyhujk")

func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		utils.Log(slog.LevelError, "❌Error", "unable to validate token", "data", gin.H{
			"error": err,
		})
		return nil, utils.CapitalizeError("invalid token")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		utils.Log(slog.LevelError, "❌Error", "invalid token claims, unable to generate token")
		return nil, utils.CapitalizeError("invalid token claims")
	}
	return claims, nil
}

func GenerateJWT(userid string) (string, time.Time, error) {
	// Set the expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims (payload) of the token
	claims := &jwt.RegisteredClaims{
		Subject:   userid,
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	// Return the signed token and the expiration time
	return signedToken, expirationTime, nil
}
