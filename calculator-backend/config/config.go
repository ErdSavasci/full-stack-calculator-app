package config

import (
	"log"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Services struct {
		Addition struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"addition"`
		Subtraction struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"subtraction"`
		Multiplication struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"multiplication"`
		Division struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"division"`
		Exponentiation struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"exponentiation"`
		SquareRoot struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"squareroot"`
		Percentage struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"percentage"`
		History struct {
			Port string `mapstructure:"port"`
		} `mapstructure:"history"`
	} `mapstructure:"services"`
}

func LoadConfig() (config Config) {
	slog.Info("Config is loading...")

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found: %v", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	slog.Info("Config loaded.")

	return
}
