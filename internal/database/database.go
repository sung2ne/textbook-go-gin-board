
func Init(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    // ...

    // 자동 마이그레이션
    if err := db.AutoMigrate(&domain.Post{}, &domain.Comment{}); err != nil {
        return nil, err
    }

    // ...
}
