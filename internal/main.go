package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/yourtion/domain-proxy/internal/base/config"
	"github.com/yourtion/domain-proxy/internal/base/define"
	"github.com/yourtion/domain-proxy/internal/base/logger"
	"github.com/yourtion/domain-proxy/internal/proxy"
)

var log *logger.Entry

func init() {
	log = logger.NewModuleLogger("main")
}

func startPProf() {
	log.Warnf("start PProf web interface on %s", config.Config.Server.PProf)
	log.Infof("Open on http://%s/debug/pprof/", config.Config.Server.PProf)
	go func() {
		if err := http.ListenAndServe(config.Config.Server.PProf, nil); err != nil {
			log.Warnln(err)
		}
	}()
}

func startServer() {
	rp := proxy.NewReverseProxyPool()
	var err error
	conf := config.Config.Server
	// 监听端口
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		log.Fatalf("listener failed: %s", err)
	}
	if config.Config.Server.Https {
		err = http.ServeTLS(listener, rp, conf.HttpsPem, conf.HttpsKey)
	} else {
		err = http.Serve(listener, rp)
	}
	if err != nil {
		if err == http.ErrServerClosed {
			log.Warnln(err)
		} else {
			log.Fatalf("listen http failed: %s", err)
		}
	}
}

func main() {
	log.Infof("pid: %d, gid: %d, uid: %d", os.Getpid(), os.Getgid(), os.Getuid())
	// 根据运行目录获取配置文件名
	configFile := "config.toml"
	workingDir := "./"

	// 载入配置
	config.LoadConfig(workingDir, configFile)
	log.Infof("server name is %s -> %s", config.Config.Server.Name, define.ServiceName)
	log.Infof("config: %+v", config.Config)

	// 切换到指定的工作目录
	if err := os.Chdir(config.Config.CWD); err != nil {
		log.Errorf("change working directory to %s failed: %s", config.Config.CWD, err)
	} else {
		log.Infof("current working directory is %s", config.Config.CWD)
	}

	// 初始化日志记录器
	logger.InitLogger(config.Config.Log.Level)

	// 判断是否需要启动 PProf
	if len(config.Config.Server.PProf) > 0 {
		startPProf()
	}

	log.Printf("Server listen on %s", config.Config.Server.Listen)
	startServer()
}
