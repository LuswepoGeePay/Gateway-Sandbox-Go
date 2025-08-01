package authservices

import (
	"fmt"
	"log/slog"
	"net/http"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	"pg_sandbox/utils"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginUser(req *auth.LoginRequest) (*auth.AuthResponse, int, error) {
	var user models.User

	// First check if user exists
	result := config.DB.Preload("Role.Permissions").Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		utils.Log(slog.LevelError, "❌Error", "Unable to Login", "detail", result.Error.Error(), "data", gin.H{
			"email": req.Email,
		})
		if result.Error == gorm.ErrRecordNotFound {
			utils.Log(slog.LevelError, "❌Error", "Unable to Login, user not found", "data", gin.H{
				"email": req.Email,
			})

			return nil, http.StatusUnauthorized, utils.CapitalizeError("Invalid email or password")
		}
		return nil, http.StatusInternalServerError, utils.CapitalizeError(result.Error.Error())
	}

	// Separate password check to avoid timing attacks
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, http.StatusInternalServerError, utils.CapitalizeError("invalid credentials")
	}

	token, tokenExpiry, err := GenerateJWT(user.ID.String())
	if err != nil {
		utils.Log(slog.LevelError, "❌Error", "Unable to Login, unable to create token", "data", gin.H{
			"email":   req.Email,
			"user_id": user.ID,
		})

		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create token: %v", err)
	}

	requestedPortal := req.Portal

	userRole := user.Role.Name

	allowedPortals := map[string][]string{
		"admin":    {"admin"},
		"merchant": {"merchant"},
	}

	if allowedPortals[userRole] == nil {
		return nil, http.StatusUnauthorized, utils.CapitalizeError(fmt.Sprintf("User role '%s' is not configured for any portal", userRole))
	}

	isAllowed := false
	for _, allowedPortal := range allowedPortals[userRole] {
		if requestedPortal == allowedPortal {
			isAllowed = true
			break
		}
	}

	var permissions []string
	for _, permission := range user.Role.Permissions {
		permissions = append(permissions, permission.Name)
	}

	if !isAllowed {
		utils.Log(slog.LevelError, "❌Error", fmt.Sprintf("User with role '%s' cannot access the '%s' portal", userRole, requestedPortal), "data", gin.H{
			"email": req.Email,
		})

		return nil, http.StatusForbidden, utils.CapitalizeError(fmt.Sprintf("User with role '%s' cannot access the '%s' portal", userRole, requestedPortal))
	}

	return &auth.AuthResponse{
		Success:       true,
		Status:        "success",
		Message:       "Login successful",
		Token:         token,
		Id:            user.ID.String(),
		Role:          user.Role.Name,
		Permissions:   permissions,
		TokenExpiry:   tokenExpiry.Format(time.RFC3339),
		EmailVerified: user.EmailVerified,
		AccountStatus: user.Status,
		Email:         user.Email,
	}, 200, nil
}
