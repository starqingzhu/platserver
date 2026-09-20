package gm

import (
	"net/http"

	"platserver/internal/dispatch"

	libdispatch "golib/dispatch"
	commonmsg "pbcommon/gen/common/msg"
	pb "pbcommon/gen/ss/msg"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
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

	payload, err := proto.Marshal(&pb.PBS2SRobotNoticeRequest{
		UserId:  req.UserId,
		BizType: req.BizType,
		Body:    req.Body,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	meta := libdispatch.Meta{Protocol: libdispatch.ProtocolHTTP, RemoteIP: ctx.ClientIP()}
	respBytes, err := dispatch.Main.Dispatch(ctx.Request.Context(), meta, dispatch.CmdRobotNoticePlayer, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	resp := &pb.PBS2RobotNoticeResponse{}
	if err := proto.Unmarshal(respBytes, resp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	if resp.MsgCode != commonmsg.MsgCode_CODE_OK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "send failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
