package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"domain-proxy/src/base/config"
	"domain-proxy/src/base/define"
	"domain-proxy/src/base/logger"
)

var log *logger.Entry
var server *http.Server

func init() {
	log = logger.NewModuleLogger("main").WithField("version", define.Version)
}

func startServer() {
	rp := NewReverseProxyPool()
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
	log.Infof("server name is %s", config.Config.Server.Name)
	log.Infof("config: %+v", config.Config)

	// 切换到指定的工作目录
	if err := os.Chdir(config.Config.CWD); err != nil {
		log.Errorf("change working directory to %s failed: %s", config.Config.CWD, err)
	} else {
		log.Infof("current working directory is %s", config.Config.CWD)
	}

	// 初始化日志记录器
	if level, err := logrus.ParseLevel(config.Config.Log.Level); err != nil {
		log.Errorf("invalid log level: %s", config.Config.Log.Level)
	} else {
		logger.Logger.SetLevel(level)
		if level >= logrus.DebugLevel {
			// 调试模式下输出带颜色的日志，方便阅读
			logger.Logger.SetFormatter(&logrus.TextFormatter{
				ForceColors: true,
			})
		}
		log.Infof("log level is %s", config.Config.Log.Level)
	}

	// 判断是否需要启动 pprof
	if len(config.Config.Server.PProf) > 0 {
		log.Warnf("start pprof web interface on %s", config.Config.Server.PProf)
		log.Infof("Open on http://%s/debug/pprof/", config.Config.Server.PProf)
		go func() {
			if err := http.ListenAndServe(config.Config.Server.PProf, nil); err != nil {
				log.Warnln(err)
			}
		}()
	}
	log.Printf("server listen on %s", config.Config.Server.Listen)

	startServer()
}
