package plat

import (
	"context"
	"fmt"
	"golib/zaplog"
	"platserver/internal/module/notify"

	commonmsg "pbcommon/gen/common/msg"
	pb "pbcommon/gen/ss/msg"
)

// ServerHandler implements pb.PlatServerServer.
type ServerHandler struct {
	pb.UnimplementedPlatServerServer
}

var globalServerHandler *ServerHandler

func InitServerHandler() *ServerHandler {
	if globalServerHandler != nil {
		return globalServerHandler
	}
	globalServerHandler = &ServerHandler{}
	zaplog.LoggerSugar.Infof("[PlatServer] rpc handler initialised")
	return globalServerHandler
}

func (h *ServerHandler) Close() error {
	zaplog.LoggerSugar.Infof("[PlatServer] rpc handler closed")
	return nil
}

func (h *ServerHandler) S2SPlayer2RobotNotice(ctx context.Context, req *pb.PBS2SRobotNoticeRequest) (*pb.PBS2RobotNoticeResponse, error) {
	zaplog.LoggerSugar.Infof("[PlatServer] S2SPlayer2RobotNotice userId=%d bizType=%s", req.UserId, req.BizType)

	userIdStr := fmt.Sprint(req.UserId)
	if err := notify.Send(userIdStr, req.BizType, req.Body); err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SPlayer2RobotNotice send failed: %v", err)
		return &pb.PBS2RobotNoticeResponse{
			MsgCode: commonmsg.MsgCode_CODE_ERROR,
		}, nil
	}

	return &pb.PBS2RobotNoticeResponse{
		MsgCode: commonmsg.MsgCode_CODE_OK,
	}, nil
}

func (h *ServerHandler) S2SGmRobotNotice(ctx context.Context, req *pb.PBS2SGMRobotNoticeRequest) (*pb.PBS2GMRobotNoticeResponse, error) {
	zaplog.LoggerSugar.Infof("[PlatServer] S2SGmRobotNotice userId=%s bizType=%s", req.Name, req.BizType)

	userIdStr := req.Name
	if err := notify.Send(userIdStr, req.BizType, req.Body); err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SGmRobotNotice send failed: %v", err)
		return &pb.PBS2GMRobotNoticeResponse{
			MsgCode: commonmsg.MsgCode_CODE_ERROR,
		}, nil
	}
	return &pb.PBS2GMRobotNoticeResponse{
		MsgCode: commonmsg.MsgCode_CODE_OK,
	}, nil
}
