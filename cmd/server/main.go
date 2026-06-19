func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid config: %v", err)
    }

    // 서버 시작
}
