package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"scootin-aboot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware creates a custom recovery middleware with structured logging
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			// Log the panic with structured logging
			utils.Error("Panic recovered",
				zap.String("error", err),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("stack", string(debug.Stack())),
			)
		}

		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := recovered.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
					strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}

		// If it's a broken pipe, we don't need to log the stack trace
		if brokenPipe {
			utils.Warn("Broken pipe error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
			)
		} else {
			// Log the full stack trace for other errors
			utils.Error("Panic recovered with stack trace",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.String("stack", string(debug.Stack())),
			)
		}

		// If the connection is dead, we can't write a status to it.
		if brokenPipe {
			c.Error(recovered.(error)) // nolint: errcheck
			c.Abort()
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
}

// DumpRequestMiddleware logs the full HTTP request for debugging
func DumpRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only dump in debug mode
		if gin.Mode() == gin.DebugMode {
			requestDump, err := httputil.DumpRequest(c.Request, true)
			if err != nil {
				utils.Error("Failed to dump request", zap.Error(err))
			} else {
				utils.Debug("HTTP Request Dump",
					zap.String("request", string(requestDump)),
				)
			}
		}
		c.Next()
	}
}
