package configs

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	OAuth      OAuthConfig      `mapstructure:"oauth"`
}
