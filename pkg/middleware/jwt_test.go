package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danisasmita/customer-search/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generate a valid token
	validToken, err := utils.GenerateJWT(123)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupHeaders   func(req *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "No Authorization Header",
			setupHeaders: func(req *http.Request) {
				// Tidak ada header yang disetel untuk menguji kasus tanpa Authorization header
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"unauthorized\"}",
		},
		{
			name: "Invalid Token",
			setupHeaders: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalidtoken")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"error\":\"unauthorized\"}",
		},
		{

			name: "Valid Token",
			setupHeaders: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer "+validToken)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success", // Ini bukan JSON, jadi pakai assert.Equal, bukan assert.JSONEq

		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(JWTAuth())
			r.GET("/protected", func(c *gin.Context) {
				c.String(http.StatusOK, "success")
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			w := httptest.NewRecorder()
			tt.setupHeaders(req)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String()) // Untuk valid token, cek string biasa
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String()) // Untuk error response, cek sebagai JSON
			}

		})
	}
}
