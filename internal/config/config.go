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
	AccessExpiry  int    `mapstructure:"access_expiry"`  // 분 단위
	RefreshExpiry int    `mapstructure:"refresh_expiry"` // 시간 단위
}

type PaginationConfig struct {
	DefaultSize int `mapstructure:"default_size"`
	MaxSize     int `mapstructure:"max_size"`
}

// DSN 데이터베이스 연결 문자열 생성
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

var cfg *Config

// Load 설정 파일 로드
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 환경 변수 연동
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GOBOARD")

	// 환경 변수 키 매핑
	viper.BindEnv("database.password", "GOBOARD_DB_PASSWORD")
	viper.BindEnv("database.host", "GOBOARD_DB_HOST")
	viper.BindEnv("jwt.secret_key", "GOBOARD_JWT_SECRET")
	viper.BindEnv("redis.addr", "GOBOARD_REDIS_ADDR")

	// 설정 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	// 구조체로 언마샬
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("설정 파싱 실패: %w", err)
	}

	log.Printf("설정 로드 완료: %s", path)
	return cfg, nil
}

// Get 전역 설정 반환
func Get() *Config {
	return cfg
}
