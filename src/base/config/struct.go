package config

type MainConfig struct {
	CWD    string       `toml:"cwd"`
	Server ServerConfig `toml:"server"`
	Log    LogConfig    `toml:"log"`
}

type ServerConfig struct {
	Name     string `toml:"name"`
	Listen   string `toml:"listen"`
	PProf    string `toml:"pprof"`
	Https    bool   `toml:"https"`
	HttpsPem string `toml:"https_pem"`
	HttpsKey string `toml:"https_key"`
}

type LogConfig struct {
	Level string `toml:"level"`
}
