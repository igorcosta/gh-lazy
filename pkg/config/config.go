package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Repo      string
	TasksFile string
	TokenFile string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("repo", "")
	viper.SetDefault("tasks_file", "")
	viper.SetDefault("token_file", ".token")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
