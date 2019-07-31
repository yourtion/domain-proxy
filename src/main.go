package main

import (
	"net/http"
	"os"

	"domain-proxy/src/base/config"
	"domain-proxy/src/base/define"
	"domain-proxy/src/base/logger"
	"domain-proxy/src/proxy"
)

var log *logger.Entry

func init() {
	log = logger.NewModuleLogger("main").WithField("version", define.Version)
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
	if err := http.ListenAndServe(config.Config.Server.Listen, rp); err != nil {
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
