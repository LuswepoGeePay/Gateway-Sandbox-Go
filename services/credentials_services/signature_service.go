package credentialsservices

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log/slog"
	"pg_sandbox/config"
	"pg_sandbox/models"
	"pg_sandbox/proto/api"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func GenerateOAuthSignature(req *api.GenerateOAuthSignatureRequest) (*string, error) {

	if req.UserId == "" {
		return nil, utils.CapitalizeError("User ID is required")
	}

	if req.ClientId == "" {
		return nil, utils.CapitalizeError("client ID is required")
	}

	if req.SecretKey == "" {
		return nil, utils.CapitalizeError("secret Key is required")
	}

	if req.ClientPin == "" {
		return nil, utils.CapitalizeError("client PIN is required")
	}

	clientString := req.ClientId + ":" + req.ClientPin

	clientSecret := req.SecretKey

	message := clientSecret + clientString

	h := hmac.New(sha256.New, []byte(clientSecret))

	h.Write([]byte(message))

	signatureBytes := h.Sum(nil)

	oauthSignature := base64.StdEncoding.EncodeToString(signatureBytes)

	tx := config.DB.Begin()

	result := tx.Where("user_id = ?", req.UserId).Model(&models.ApiKeys{}).Update("o_auth_signature", oauthSignature)

	if result.Error != nil {
		tx.Rollback()
		utils.Log(slog.LevelError, "❌Error", "unable to add oauth signature to db ", "data", gin.H{
			"error":        result.Error,
			"request_body": req,
		})
		return nil, utils.CapitalizeError("Error adding signature to profile")
	}

	tx.Commit()

	return &oauthSignature, nil
}
