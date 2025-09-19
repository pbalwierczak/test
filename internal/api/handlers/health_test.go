package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.HealthCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	expectedResponse := `{"service":"scootin-aboot","status":"healthy"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
