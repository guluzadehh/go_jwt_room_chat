package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	JWT         JWTCfg     `yaml:"jwt"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	Redis       RedisCfg   `yaml:"redis"`
	Chat        Chat       `yaml:"chat"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWTCfg struct {
	SecretKey string     `yaml:"-" env:"JWT_SECRET_KEY" env-required:"true"`
	Access    AccessCfg  `yaml:"access"`
	Refresh   RefreshCfg `yaml:"refresh"`
}

type AccessCfg struct {
	Expire time.Duration `yaml:"expire" env-default:"1h"`
}

type RefreshCfg struct {
	EncryptSecretKey string        `yaml:"-" env:"JWT_REFRESH_ENCRYPT_SECRET_KEY" env-required:"true"`
	Expire           time.Duration `yaml:"expire" env-default:"168h"`
	CookieName       string        `yaml:"cookie_name" env-default:"jwt_refresh"`
}

type RedisCfg struct {
	Address   string `yaml:"address" env-default:"localhost:6379"`
	Password  string `yaml:"-" env:"REDIS_PASSWORD" env-required:"true"`
	DefaultDB int    `yaml:"default_db" env-default:"0"`
}

type Chat struct {
	Room       RoomCfg       `yaml:"room"`
	PongWait   time.Duration `yaml:"pong_wait" env-default:"60s"`
	PingPeriod time.Duration `yaml:"ping_period" env-default:"54s"`
	WriteWait  time.Duration `yaml:"write_wait" env-default:"10s"`
}

type RoomCfg struct {
	Capacity int `yaml:"capacity" env-default:"16"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file `%s` does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config file `%s` and env variables\n\t%s", configPath, err)
	}

	return &cfg
}
