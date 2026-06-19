type HealthStatus struct {
    Status     string                 `json:"status"`
    Checks     map[string]CheckResult `json:"checks"`
    Version    string                 `json:"version,omitempty"`
    Uptime     string                 `json:"uptime,omitempty"`
}

type CheckResult struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
    Latency string `json:"latency,omitempty"`
}

var startTime = time.Now()

func (h *HealthHandler) Detailed(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    checks := make(map[string]CheckResult)
    overallStatus := "ok"

    // 데이터베이스 체크
    dbCheck := h.checkDatabase(ctx)
    checks["database"] = dbCheck
    if dbCheck.Status != "ok" {
        overallStatus = "degraded"
    }

    // Redis 체크 (있다면)
    // redisCheck := h.checkRedis(ctx)
    // checks["redis"] = redisCheck

    status := HealthStatus{
        Status:  overallStatus,
        Checks:  checks,
        Version: "1.0.0",
        Uptime:  time.Since(startTime).String(),
    }

    statusCode := http.StatusOK
    if overallStatus != "ok" {
        statusCode = http.StatusServiceUnavailable
    }

    c.JSON(statusCode, status)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) CheckResult {
    start := time.Now()

    sqlDB, err := h.db.DB()
    if err != nil {
        return CheckResult{
            Status:  "error",
            Message: err.Error(),
        }
    }

    if err := sqlDB.PingContext(ctx); err != nil {
        return CheckResult{
            Status:  "error",
            Message: err.Error(),
        }
    }

    return CheckResult{
        Status:  "ok",
        Latency: time.Since(start).String(),
    }
}
