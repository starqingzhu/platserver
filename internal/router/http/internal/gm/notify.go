package gm

import (
	"net/http"
	"platserver/internal/module/notify"

	"github.com/gin-gonic/gin"
)

type robotNoticeRequest struct {
	UserId  int64  `json:"userId" binding:"required"`
	BizType string `json:"bizType" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func RobotNotice(ctx *gin.Context) {
	var req robotNoticeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "invalid params: " + err.Error()})
		return
	}

	if err := notify.Send(req.UserId, req.BizType, req.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
