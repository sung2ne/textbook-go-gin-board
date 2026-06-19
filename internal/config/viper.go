package config

import (
    "github.com/spf13/viper"
)

func LoadWithViper() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    // 환경 변수 자동 읽기
    viper.AutomaticEnv()

    // 설정 파일 읽기 (없어도 에러 아님)
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    cfg := &Config{}
    if err := viper.Unmarshal(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}
