type App struct {
    server   *http.Server
    db       *gorm.DB
    cleanups []func() error
}

func NewApp(cfg *config.Config) (*App, error) {
    app := &App{}

    // 데이터베이스
    db, err := database.Connect(cfg.Database.URL)
    if err != nil {
        return nil, err
    }
    app.db = db
    app.cleanups = append(app.cleanups, func() error {
        sqlDB, _ := db.DB()
        return sqlDB.Close()
    })

    // 서버
    router := setupRouter(cfg, db)
    app.server = &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
        Handler: router,
    }

    return app, nil
}

func (a *App) Run() error {
    return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
    // 서버 종료
    if err := a.server.Shutdown(ctx); err != nil {
        return err
    }

    // 역순으로 정리 (나중에 초기화된 것 먼저)
    for i := len(a.cleanups) - 1; i >= 0; i-- {
        if err := a.cleanups[i](); err != nil {
            slog.Error("cleanup failed", "error", err)
        }
    }

    return nil
}
