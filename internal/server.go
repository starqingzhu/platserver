package internal

import (
	"common/config"
	"common/defines"
	cetcd "common/etcd"
	"fmt"
	gohttp "golib/http"
	"golib/node"
	"golib/utils"
	"golib/yamlcfg"
	"golib/zaplog"
	"path/filepath"
	"platserver/internal/module/notify"
	hrouter "platserver/internal/router/http"
	rpcservice "platserver/internal/router/rpc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	GitCommitSha1 string
	CompileDate   string
}

func (s *Server) Name() string {
	return defines.PlatServer.String()
}

func (s *Server) Run(closeChan chan struct{}) {
	gohttp.Main.Run()
	<-closeChan
}

func (s *Server) LoadConfig() {
	cfg, ok := yamlcfg.LoadYamlCfg()
	if !ok {
		zaplog.LoggerSugar.Fatalf("LoadConfig: failed to load YAML configuration file")
		return
	}

	config.Default.Game = cfg.Game
	config.Default.Cluster = cfg.Cluster
	config.Default.ModuleName = cfg.Module
	config.Default.ConfigDir = cfg.ConfigDir
	config.Default.HttpCfg = cfg.HttpCfg
	config.Default.RpcCfg = cfg.RpcCfg
	config.Default.LogCfg = cfg.LogCfg
	config.Default.PprofCfg = cfg.PprofCfg

	if len(cfg.Etcd) > 0 {
		etcdCfg := cfg.Etcd[0]
		etcdHost := fmt.Sprintf("%s:%d", etcdCfg.Host, etcdCfg.Port)
		config.Default.EtcdCfg.EndPoints = []string{etcdHost}
		config.Default.EtcdCfg.Username = etcdCfg.Username
		config.Default.EtcdCfg.Password = etcdCfg.Password
	}

	zaplog.LoggerSugar.Infof("LoadConfig: configuration loaded successfully")
}

func (s *Server) OnInit() {
	s.initServerNodeInfo()

	cetcd.InitAndWatchServerType(false, defines.PlatServer.String())

	s.initConfigFiles()

	rpcservice.InitAllRPC()

	s.initHTTPServer()

	cetcd.RefreshNodeStateInNormal()

	zaplog.LoggerSugar.Infof("platserver init complete, commitId:%s, compileDate:%s", s.GitCommitSha1, s.CompileDate)
}

func (s *Server) OnClose() {
	notify.Shutdown()
	rpcservice.Close()
	cetcd.Close()
}

// initConfigFiles 加载 config 目录下的业务配置文件。
func (s *Server) initConfigFiles() {
	execPath, ok := utils.GetExecCurrentPath()
	if !ok {
		zaplog.LoggerSugar.Fatalf("initConfigFiles: failed to get executable path")
	}

	relativeConfigDir := config.Default.GetConfigDir()
	configDir, err := filepath.Abs(filepath.Join(execPath, relativeConfigDir))
	if err != nil {
		zaplog.LoggerSugar.Fatalf("initConfigFiles: failed to resolve config path, err: %v", err)
	}

	zaplog.LoggerSugar.Infof("initConfigFiles: loading from %s", configDir)

	if err := notify.LoadConfig(configDir); err != nil {
		zaplog.LoggerSugar.Fatalf("initConfigFiles: failed to load notify config, err: %v", err)
	}
}

func (s *Server) initHTTPServer() {
	gin.SetMode(gin.ReleaseMode)
	gohttp.InitMainServer(config.Default.HttpCfg)
	hrouter.InitRouter(gohttp.Main.Engine)
	zaplog.LoggerSugar.Infof("initHTTPServer: HTTP server initialized on %s", config.Default.HttpCfg.HttpListenAddr)
}

func (s *Server) initServerNodeInfo() {
	rpcAddr, err := utils.GenerateRegServerAddr(config.Default.RpcCfg.Server.Address)
	if err != nil {
		zaplog.LoggerSugar.Fatalf("initServerNodeInfo: failed to generate rpc addr, err: %v", err)
		panic(err)
	}

	node.Init(
		config.Default.Cluster,
		rpcAddr,
		rpcAddr,
		defines.PlatServer.String(),
		node.ServiceState_Normal,
		"",
		s.GitCommitSha1,
		s.CompileDate,
	)

	zaplog.LoggerSugar.Infof("initServerNodeInfo: node initialized, cluster=%s, addr=%s, type=%s",
		config.Default.Cluster, rpcAddr, defines.PlatServer.String())
}
