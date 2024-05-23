package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ygpkg/yg-go/apis/constants"
	"github.com/ygpkg/yg-go/apis/runtime"
	"github.com/ygpkg/yg-go/logs"
)

// Logger .
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqid := ctx.GetString(constants.CtxKeyRequestID)
		logger := logs.RequestLogger(reqid)

		ctx.Set(constants.CtxKeyLogger, logger)
		currReq := fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())
		for _, whitelistItem := range []string{
			"POST /apis/zmdevice/v1/ping",
			"POST /apis/zmdevice/v1/task_logs",
			"POST /apis/p/zmrobot.KeepTaskHeartbeat",
		} {
			if whitelistItem == currReq {
				ctx.Next()
				return
			}
		}
		start := time.Now()
		ctx.Next()
		cost := time.Since(start)

		logger.Infow(fmt.Sprint(ctx.Writer.Status()),
			"method", ctx.Request.Method,
			"uri", ctx.Request.RequestURI,
			"reqsize", ctx.Request.ContentLength,
			"latency", fmt.Sprintf("%.3f", cost.Seconds()),
			"clientip", runtime.GetRealIP(ctx.Request),
			"respsize", ctx.Writer.Size(),
		)
	}
}
