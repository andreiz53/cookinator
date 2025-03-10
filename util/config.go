package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all environment variables of the application
type Config struct {
	DBSource          string        `mapstructure:"DB_SOURCE"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DURATION"`
}

// LoadConfig reads the configuration file using viper
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	return config, err
}
