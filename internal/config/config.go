package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	JWT    JWTConfig

	PasswordSalt string
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	SigningKey      string
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string
	DBName   string
	SSLMode  string `mapstructure:"sslmode"`
	Password string
}

func Init(configPath string) (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := unmarshal(&config); err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	config.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	config.DB.Username = os.Getenv("POSTGRES_USER")
	config.DB.DBName = os.Getenv("POSTGRES_DB")
	config.PasswordSalt = os.Getenv("password_salt")
	config.JWT.SigningKey = os.Getenv("signing_key")

	return &config, nil
}

func unmarshal(config *Config) error {
	if err := viper.UnmarshalKey("server", &config.Server); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("db", &config.DB); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &config.JWT); err != nil {
		return err
	}

	return nil
}
