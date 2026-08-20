package users

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func validBearer() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   "debug-user",
		"client_id": "debug-client",
		"exp":       time.Now().UTC().Add(time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte("SuperSecretKeyForARobustSystem"))
	return "Bearer " + s
}

func TestNameLookUpReproduction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/mobile-money/name-lookup/:phone", NameLookUpHandler)
	r.GET("/v1/mobile-money/name-lookup/:number", NameLookUpHandler)

	auth := validBearer()
	cases := []struct {
		name    string
		method  string
		path    string
		headers map[string]string
	}{
		{"no-auth", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424", nil},
		{"spring-typical-get", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424", map[string]string{
			"Authorization": auth,
			"Accept":        "application/json, application/*+json",
		}},
		{"bearer-and-accept-json", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424", map[string]string{
			"Authorization": auth,
			"Accept":        "application/json",
		}},
		{"charset-content-type", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424", map[string]string{
			"Authorization": auth,
			"Accept":        "application/json",
			"Content-Type":  "application/json; charset=utf-8",
		}},
		{"full-exact-headers", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424", map[string]string{
			"Authorization": auth,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		}},
		{"trailing-quote", http.MethodGet, "/api/v1/mobile-money/name-lookup/260765631424'", map[string]string{
			"Authorization": auth,
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			t.Logf("%s status=%d body=%s", tc.name, w.Code, w.Body.String())
		})
	}
}
