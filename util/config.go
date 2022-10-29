package util

import "github.com/spf13/viper"

// Config contains all the configuration for the application
// The values are read by viper from a config file or environment variables
type Config struct {
	// Address is the address to listen on
	ServerAddress string `mapstructure:"SERVER_ADDRESS" json:"server_address" yaml:"server_address"`
	// DatabaseDriver is the database driver to use
	DatabaseDriver string `mapstructure:"DATABASE_DRIVER" json:"database_driver" yaml:"database_driver"`
	// DatabaseSource is the database connection string
	DatabaseSource string `mapstructure:"DATABASE_SOURCE" json:"database_source" yaml:"database_source"`
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
