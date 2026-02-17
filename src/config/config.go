package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Password PasswordConfig
	Cors     CorsConfig
	Logger   LoggerConfig
	JWT      JWTConfig
	Pdf      PdfConfig
}

type ServerConfig struct {
	InternalPort string
	ExternalPort string
	RunMode      string
	Domain       string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type CorsConfig struct {
	AllowOrigins string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

type JWTConfig struct {
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
	Secret                     string
	RefreshSecret              string
}

type PdfConfig struct {
	Url string
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := LoadConfig(cfgPath, "yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg, err := ParseConfig(v)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		cfg.Server.ExternalPort = envPort
		log.Printf("Set external port from environment -> %s", cfg.Server.ExternalPort)
	}
	return cfg
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/config/config-docker"
	} else if env == "production" {
		return "/config/config-production"
	} else {
		return "./config/config-development"
	}
}
func LoadConfig(fileName string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.AddConfigPath(filepath.Dir(fileName))
	v.SetConfigName(filepath.Base(fileName))
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}
