package middleware

import (
	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a custom logging middleware using zap logger
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log using zap with structured logging
		utils.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Time("timestamp", param.TimeStamp),
		)
		return ""
	})
}
