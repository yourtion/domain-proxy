package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var log *Entry

func init() {
	Logger = logrus.New()
	// 使用 JSON 格式记录
	// Logger.SetFormatter(&logrus.JSONFormatter{})
	// 输出到 stdout
	Logger.SetOutput(os.Stdout)

	log = NewModuleLogger("logger")
}

type Entry = logrus.Entry

func WithFields(key string, value interface{}) *Entry {
	return Logger.WithField(key, value)
}

func NewModuleLogger(name string) *Entry {
	return Logger.WithField("module", name)
}

func InitLogger(logLevel string) {
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		log.Errorf("invalid log level: %s", logLevel)
	} else {
		Logger.SetLevel(level)
		if level >= logrus.DebugLevel {
			// 调试模式下输出带颜色的日志，方便阅读
			Logger.SetFormatter(&logrus.TextFormatter{
				ForceColors: true,
			})
		}
		log.Infof("log level is %s", logLevel)
	}
}
