package config

import (
	"github.com/BurntSushi/toml"

	"github.com/yourtion/domain-proxy/internal/base/logger"
)

var log *logger.Entry
var Config MainConfig

func init() {
	log = logger.NewModuleLogger("config")
}

func LoadConfig(workingDir string, file string) {
	log.Infof("load config from file: %s", file)
	_, err := toml.DecodeFile(file, &Config)
	if err != nil {
		log.Fatalf("load config failed: %s", err)
	}

	if len(Config.CWD) < 1 {
		Config.CWD = workingDir
	}
	Config.Loaded = true
}

func MockConfig(conf MainConfig) {
	if Config.Loaded {
		panic("Mock loaded config")
	}
	log.Infof("mock config: %v", conf)
	Config = conf
}
