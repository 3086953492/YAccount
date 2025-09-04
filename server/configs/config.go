package configs

import "github.com/3086953492/YaBase/configs"

type Config struct {
	Server     configs.ServerConfig     `mapstructure:"server"`
	Database   configs.DatabaseConfig   `mapstructure:"database"`
	Redis      configs.RedisConfig      `mapstructure:"redis"`
	JWT        configs.JWTConfig        `mapstructure:"jwt"`
	Log        configs.LogConfig        `mapstructure:"log"`
	Middleware configs.MiddlewareConfig `mapstructure:"middleware"`
	OAuth      configs.OAuthConfig      `mapstructure:"oauth"`
}
