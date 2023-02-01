package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config contains all the configuration for the application
// The values are read by viper from a config file or environment variables
type Config struct {
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	DatabaseDriver       string        `mapstructure:"DATABASE_DRIVER"`
	DatabaseSource       string        `mapstructure:"DATABASE_SOURCE"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
}

// LoadConfig loads the configuration from a config file or environment variables
func LoadConfig(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	var cfg Config
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
