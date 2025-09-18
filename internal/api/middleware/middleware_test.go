package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandlerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlerMiddleware())

	router.GET("/test-error", func(c *gin.Context) {
		c.Error(NewAPIError(http.StatusBadRequest, "Test error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-error", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Test error")
}

func TestValidateJSONMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateJSON())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test with valid JSON content type
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with invalid content type
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "text/plain")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
