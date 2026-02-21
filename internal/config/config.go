package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Pagination PaginationConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	SecretKey     string `mapstructure:"secret_key"`
	AccessExpiry  int    `mapstructure:"access_expiry"`
	RefreshExpiry int    `mapstructure:"refresh_expiry"`
}

type PaginationConfig struct {
	DefaultSize int `mapstructure:"default_size"`
	MaxSize     int `mapstructure:"max_size"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

var cfg *Config

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GOBOARD")

	viper.BindEnv("database.password", "GOBOARD_DB_PASSWORD")
	viper.BindEnv("database.host", "GOBOARD_DB_HOST")
	viper.BindEnv("jwt.secret_key", "GOBOARD_JWT_SECRET")
	viper.BindEnv("redis.addr", "GOBOARD_REDIS_ADDR")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("설정 파싱 실패: %w", err)
	}

	log.Printf("설정 로드 완료: %s", path)
	return cfg, nil
}

func Get() *Config {
	return cfg
}
