package cache

func InitRedisWithPool(addr, password string) error {
    redisClient = redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     password,
        DB:           0,
        PoolSize:     100,              // 최대 연결 수
        MinIdleConns: 10,               // 최소 유휴 연결
        PoolTimeout:  4 * time.Second,  // 연결 대기 시간
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return redisClient.Ping(ctx).Err()
}
