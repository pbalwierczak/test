package middleware

import (
	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		utils.Info("HTTP Request",
			utils.String("method", param.Method),
			utils.String("path", param.Path),
			utils.Int("status", param.StatusCode),
			utils.Duration("latency", param.Latency),
			utils.String("client_ip", param.ClientIP),
			utils.String("user_agent", param.Request.UserAgent()),
			utils.Time("timestamp", param.TimeStamp),
		)
		return ""
	})
}
