package authservices

import (
	"fmt"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/auth"
	"pg_sandbox/utils"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("tyhsdndfuadbajsddoewkmdiedwnnpewesedrftgyhujk")

func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, utils.CapitalizeError("invalid token")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
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

func RegisterUser(req *auth.RegisterRequest) (*string, error) {

	tx := config.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, utils.CapitalizeError("unable to hash password")
	}

	var role models.Role
	result := tx.Where("name = ?", req.Role).First(&role)
	if result.Error != nil {
		return nil, utils.CapitalizeError("unable to find role.")
	}

	userId := uuid.New()
	user := models.User{
		ID:       userId,
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
		Status:   "active",
	}

	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return nil, utils.CapitalizeError(result.Error.Error())
	}

	clientID := uuid.New().String()
	clientSecret := uuid.New().String() // You can use a more secure approach to generate the secret

	apiKey := models.ApiKeys{
		ID:           uuid.New(),
		UserID:       user.ID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		IsActive:     true,
	}

	if err := tx.Create(&apiKey).Error; err != nil {
		tx.Rollback()
		return nil, utils.CapitalizeError("unable to create API keys")
	}

	tx.Commit()

	userIdStr := userId.String()

	return &userIdStr, nil
}

func LoginUser(req *auth.LoginRequest) (*auth.AuthResponse, error) {
	var user models.User

	// First check if user exists
	result := config.DB.Preload("Role.Permissions").Where("email = ?", req.Email).First(&user)
	if result.Error != nil {

		if result.Error == gorm.ErrRecordNotFound {
			return nil, utils.CapitalizeError("Invalid email or password")
		}
		return nil, utils.CapitalizeError(result.Error.Error())
	}

	// Separate password check to avoid timing attacks
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.CapitalizeError("invalid credentials")
	}

	token, tokenExpiry, err := GenerateJWT(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	var permissions []string
	for _, permission := range user.Role.Permissions {
		permissions = append(permissions, permission.Name)
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
		KycStatus:     "",
	}, nil
}

func GetUsers(req *auth.GetUsersRequest) (*auth.GetUsersResponse, error) {

	var authModel []models.User
	offset := (req.Page - 1) * req.PageSize

	tx := config.DB.Begin()

	var totalUsers int64

	query := tx.Model(&models.User{}).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ? ", "admin")

	// Count total users matching the criteria
	if err := query.Count(&totalUsers).Error; err != nil {
		utils.Log(slog.LevelError, "Failed to count users", "error", err.Error())
		tx.Rollback()
		return nil, utils.CapitalizeError(err.Error())
	}

	totalPages := int32((totalUsers + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Fetch users with pagination
	if err := query.Limit(int(req.PageSize)).Offset(int(offset)).Find(&authModel).Error; err != nil {
		utils.Log(slog.LevelError, "Failed to retrieve users", "error", err.Error())
		tx.Rollback()
		return nil, utils.CapitalizeError(err.Error())
	}

	tx.Commit()

	var pbUser []*auth.User
	for _, authI := range authModel {
		pbUser = append(pbUser, &auth.User{
			Id:       authI.ID.String(),
			Fullname: authI.Fullname,
			Email:    authI.Email,
			Phone:    authI.Phone,
			Role:     authI.Role.Name,
		})
	}

	return &auth.GetUsersResponse{
		User:        pbUser,
		TotalPages:  totalPages,
		HasMore:     req.Page < totalPages,
		CurrentPage: req.Page,
	}, nil
}
