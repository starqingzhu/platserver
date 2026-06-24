package rpc

import (
	"common/config"
	"context"
	rpccommon "golib/rpc/common"
	rpcserver "golib/rpc/server"
	"golib/zaplog"
	plathandler "platserver/internal/router/rpc/plat"
	"time"

	pb "pbcommon/gen/ss/msg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcServer      *rpcserver.Server
	platSvrHandler  *plathandler.ServerHandler
)

func InitAllRPC() {
	platSvrHandler = plathandler.InitServerHandler()
	startGRPCServer()
}

func startGRPCServer() {
	rpcCfg := &config.Default.RpcCfg
	listenAddr := rpcCfg.Server.Address
	if listenAddr == "" {
		zaplog.LoggerSugar.Errorf("[rpc] RpcCfg.Server.Address not set")
		return
	}

	grpcConfig := rpccommon.DefaultServerConfig(listenAddr)
	grpcConfig.Keepalive = keepalive.ServerParameters{
		Time:    time.Duration(rpcCfg.GetKeepaliveTime()) * time.Second,
		Timeout: time.Duration(rpcCfg.GetKeepaliveTimeout()) * time.Second,
	}
	grpcConfig.EnforcementPolicy = keepalive.EnforcementPolicy{
		MinTime:             time.Duration(rpcCfg.GetEnforcementMinTime()) * time.Second,
		PermitWithoutStream: false,
	}
	grpcConfig.EnableReflection = rpcCfg.Server.EnableReflection

	grpcServer = rpcserver.NewServer(grpcConfig)
	grpcServer.RegisterService(func(s *grpc.Server) {
		pb.RegisterPlatServerServer(s, platSvrHandler)
	})

	if err := grpcServer.StartAsync(); err != nil {
		zaplog.LoggerSugar.Errorf("[rpc] failed to start gRPC server: %v", err)
	} else {
		zaplog.LoggerSugar.Infof("[rpc] gRPC server started on %s", grpcConfig.Address)
	}
}

func Close() {
	if grpcServer == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := grpcServer.Stop(ctx); err != nil {
		zaplog.LoggerSugar.Warnf("[rpc] gRPC server stop error: %v, force stopping", err)
		grpcServer.ForceStop()
	}
	zaplog.LoggerSugar.Infof("[rpc] gRPC server stopped")

	if platSvrHandler != nil {
		if err := platSvrHandler.Close(); err != nil {
			zaplog.LoggerSugar.Warnf("[rpc] plat handler close error: %v", err)
		}
	}
}
