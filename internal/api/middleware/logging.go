package middleware

import (
	"scootin-aboot/internal/logger"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			logger.String("method", param.Method),
			logger.String("path", param.Path),
			logger.Int("status", param.StatusCode),
			logger.Duration("latency", param.Latency),
			logger.String("client_ip", param.ClientIP),
			logger.String("user_agent", param.Request.UserAgent()),
			logger.Time("timestamp", param.TimeStamp),
		)
		return ""
	})
}
