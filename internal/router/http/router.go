package http

import (
	"common/defines"
	"fmt"
	"golib/ginm"
	"golib/ginpprof"
	"golib/zaplog"
	"net/http"
	"platserver/internal/router/http/internal/gm"

	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	engine.Use(
		ginm.Recovery(),
		ginm.AccessLogWithZap(zaplog.Logger),
	)

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	})

	ginpprof.Wrap(engine)

	engine.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// gm 测试接口
	gmGroup := engine.Group(fmt.Sprintf("gm/%s", defines.PlatServer.String()))
	{
		gmGroup.POST("/robotNotice", gm.RobotNotice) // 测试机器人通知
	}
}
