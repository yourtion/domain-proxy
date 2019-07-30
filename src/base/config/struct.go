package config

type MainConfig struct {
	CWD    string       `toml:"cwd"`
	Server ServerConfig `toml:"server"`
	Log    LogConfig    `toml:"log"`
}

type ServerConfig struct {
	Name   string `toml:"name"`
	Listen string `toml:"listen"`
	URL    string `toml:"url"`
	Prefix string `toml:"prefix"`
	PProf  string `toml:"pprof"`
}

type LogConfig struct {
	Level string `toml:"level"`
}
