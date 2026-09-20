package plat

import (
	"context"
	"golib/zaplog"
	"platserver/internal/dispatch"

	libdispatch "golib/dispatch"
	commonmsg "pbcommon/gen/common/msg"
	pb "pbcommon/gen/ss/msg"

	"google.golang.org/protobuf/proto"
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

	payload, err := proto.Marshal(req)
	if err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SPlayer2RobotNotice marshal failed: userId=%d, err=%v", req.UserId, err)
		return &pb.PBS2RobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}

	meta := libdispatch.Meta{Protocol: libdispatch.ProtocolGRPC}
	respBytes, err := dispatch.Main.Dispatch(ctx, meta, dispatch.CmdRobotNoticePlayer, payload)
	if err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SPlayer2RobotNotice dispatch failed: userId=%d, err=%v", req.UserId, err)
		return &pb.PBS2RobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}

	resp := &pb.PBS2RobotNoticeResponse{}
	if err := proto.Unmarshal(respBytes, resp); err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SPlayer2RobotNotice unmarshal resp failed: userId=%d, err=%v", req.UserId, err)
		return &pb.PBS2RobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}
	return resp, nil
}

func (h *ServerHandler) S2SGmRobotNotice(ctx context.Context, req *pb.PBS2SGMRobotNoticeRequest) (*pb.PBS2GMRobotNoticeResponse, error) {
	zaplog.LoggerSugar.Infof("[PlatServer] S2SGmRobotNotice userId=%s bizType=%s", req.Name, req.BizType)

	payload, err := proto.Marshal(req)
	if err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SGmRobotNotice marshal failed: name=%s, err=%v", req.Name, err)
		return &pb.PBS2GMRobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}

	meta := libdispatch.Meta{Protocol: libdispatch.ProtocolGRPC}
	respBytes, err := dispatch.Main.Dispatch(ctx, meta, dispatch.CmdRobotNoticeGM, payload)
	if err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SGmRobotNotice dispatch failed: name=%s, err=%v", req.Name, err)
		return &pb.PBS2GMRobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}

	resp := &pb.PBS2GMRobotNoticeResponse{}
	if err := proto.Unmarshal(respBytes, resp); err != nil {
		zaplog.LoggerSugar.Errorf("[PlatServer] S2SGmRobotNotice unmarshal resp failed: name=%s, err=%v", req.Name, err)
		return &pb.PBS2GMRobotNoticeResponse{MsgCode: commonmsg.MsgCode_CODE_ERROR}, nil
	}
	return resp, nil
}
