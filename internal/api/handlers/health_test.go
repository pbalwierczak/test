package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new health handler
	handler := NewHealthHandler()

	// Create a test context with response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	handler.HealthCheck(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Assert the JSON response
	expectedResponse := `{"service":"scootin-aboot","status":"healthy"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestNewHealthHandler(t *testing.T) {
	// Test that NewHealthHandler returns a non-nil handler
	handler := NewHealthHandler()
	assert.NotNil(t, handler)
	assert.IsType(t, &HealthHandler{}, handler)
}
