package handler

import (
	"context"
	"fmt"

	"platserver/internal/dispatch"
	"platserver/internal/module/notify"

	libdispatch "golib/dispatch"
	dispatchproto "golib/dispatch/proto"
	"golib/zaplog"

	commonmsg "pbcommon/gen/common/msg"
	pb "pbcommon/gen/ss/msg"
)

// registerRobotNotice 注册机器人通知相关命令到派发器。
func registerRobotNotice(d *dispatchproto.Dispatcher) error {
	if err := dispatchproto.Register[*pb.PBS2SRobotNoticeRequest, *pb.PBS2RobotNoticeResponse](
		d, dispatch.CmdRobotNoticePlayer, handleRobotNoticePlayer,
	); err != nil {
		return err
	}
	return dispatchproto.Register[*pb.PBS2SGMRobotNoticeRequest, *pb.PBS2GMRobotNoticeResponse](
		d, dispatch.CmdRobotNoticeGM, handleRobotNoticeGM,
	)
}

func handleRobotNoticePlayer(_ context.Context, _ libdispatch.Meta, req *pb.PBS2SRobotNoticeRequest, resp *pb.PBS2RobotNoticeResponse) error {
	userIdStr := fmt.Sprint(req.UserId)
	if err := notify.Send(userIdStr, req.BizType, req.Body); err != nil {
		zaplog.LoggerSugar.Errorf("[RobotNotice] player send failed: userId=%s, bizType=%s, err=%v", userIdStr, req.BizType, err)
		resp.MsgCode = commonmsg.MsgCode_CODE_ERROR
		return nil
	}
	resp.MsgCode = commonmsg.MsgCode_CODE_OK
	return nil
}

func handleRobotNoticeGM(_ context.Context, _ libdispatch.Meta, req *pb.PBS2SGMRobotNoticeRequest, resp *pb.PBS2GMRobotNoticeResponse) error {
	if err := notify.Send(req.Name, req.BizType, req.Body); err != nil {
		zaplog.LoggerSugar.Errorf("[RobotNotice] gm send failed: name=%s, bizType=%s, err=%v", req.Name, req.BizType, err)
		resp.MsgCode = commonmsg.MsgCode_CODE_ERROR
		return nil
	}
	resp.MsgCode = commonmsg.MsgCode_CODE_OK
	return nil
}
