
func Load(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")

    // 환경 변수 바인딩
    viper.AutomaticEnv()
    viper.SetEnvPrefix("GOBOARD")  // GOBOARD_DATABASE_PASSWORD 형식으로 사용

    // 환경 변수 키 매핑
    viper.BindEnv("database.password", "GOBOARD_DB_PASSWORD")
    viper.BindEnv("database.host", "GOBOARD_DB_HOST")

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
    }

    cfg = &Config{}
    if err := viper.Unmarshal(cfg); err != nil {
        return nil, fmt.Errorf("설정 파싱 실패: %w", err)
    }

    return cfg, nil
}
