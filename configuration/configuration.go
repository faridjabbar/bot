package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	ApiKey  string `mapstructure:"API_KEY"`
}

func LoadConfig() (config Configuration, err error) {
	viper.SetConfigFile("./configuration/.env")
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
