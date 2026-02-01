package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBConn  string `mapstructure:"DB_CONN"`
	AppPort string `mapstructure:"APP_PORT"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Bind environment variables
	_ = viper.BindEnv("DB_CONN")
	_ = viper.BindEnv("APP_PORT")
	_ = viper.BindEnv("PORT") // For Render

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Override AppPort if PORT env var is set (Render)
	if port := os.Getenv("PORT"); port != "" {
		config.AppPort = port
	}

	// Default port
	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	return &config
}
