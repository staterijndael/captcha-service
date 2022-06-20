package config

import (
	"time"
)

type Config struct {
	App     `toml:"app"`
	Captcha `toml:"captcha"`
	Redis   `toml:"redis"`
}

type App struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

type Captcha struct {
	Words []string `toml:"words"`
}

type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`

	SessionTTL time.Duration `toml:"session_ttl"`
}
