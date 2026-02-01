package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBConn     string `mapstructure:"DB_CONN"`
	AppPort    string `mapstructure:"APP_PORT"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Handle Render's PORT env var
	if err := viper.BindEnv("PORT"); err != nil {
		log.Println("Failed to bind PORT env var")
	}
	
	// Explicitly bind DB keys to ensure environment variables are picked up
	viper.BindEnv("DB_CONN")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Override AppPort if PORT env var is set (Render)
	if port := viper.GetString("PORT"); port != "" {
		config.AppPort = port
	}

	return &config
}
