package commonservices

import (
	"log/slog"
	"mime"
	"net/http"
	"strings"

	tokenservices "pg_sandbox/services/token_services"
	"pg_sandbox/utils"

	"github.com/gin-gonic/gin"
)

func headerAllowsJSON(header string) bool {
	if strings.TrimSpace(header) == "" {
		return false
	}
	for _, part := range strings.Split(header, ",") {
		media, _, err := mime.ParseMediaType(strings.TrimSpace(part))
		if err != nil {
			continue
		}
		if media == "application/json" || media == "application/*" || media == "*/*" {
			return true
		}
	}
	return false
}

func CheckEssentialHeaders(c *gin.Context) {

	acceptedType := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")
	authorization := c.GetHeader("Authorization")
	method := c.Request.Method
	requiresBody := method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch

	if authorization == "" {
		// #region agent log
		utils.AgentDebugLog("common_service.go:CheckEssentialHeaders", "abort missing Authorization", "D", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
		// #endregion
		c.JSON(400, gin.H{
			"message": "unauthenticated",
		})
		c.Abort()
		return

	}

	utils.Log(slog.LevelInfo, "Authorization Header: ", authorization)

	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	err := tokenservices.ValidateOAuthToken(tokenString)
	if err != nil {

		utils.RespondWithError(c, 401, "Invalid Token")
		c.Abort()
		return
	}

	if requiresBody && !headerAllowsJSON(contentType) {
		// #region agent log
		utils.AgentDebugLog("common_service.go:CheckEssentialHeaders", "abort Content-Type mismatch", "A", map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"contentType": contentType,
			"hasCharset":  strings.Contains(contentType, "charset"),
			"isGET":       c.Request.Method == "GET",
			"runId":       "post-fix",
		})
		// #endregion
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Content-Type is application/json"},
			},
		})
		c.Abort()
		return

	}

	if acceptedType != "" && !headerAllowsJSON(acceptedType) {
		// #region agent log
		utils.AgentDebugLog("common_service.go:CheckEssentialHeaders", "abort Accept mismatch", "B", map[string]interface{}{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"acceptedType": acceptedType,
			"runId":        "post-fix",
		})
		// #endregion
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Validation failed.",
			"errors": gin.H{
				"Content-Type": []string{"Expected Accepted type is application/json"},
			},
		})
		c.Abort()
		return
	}

	// #region agent log
	utils.AgentDebugLog("common_service.go:CheckEssentialHeaders", "headers passed", "A,B,D", map[string]interface{}{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"contentType": contentType,
		"accept":      acceptedType,
		"runId":       "post-fix",
	})
	// #endregion

}
