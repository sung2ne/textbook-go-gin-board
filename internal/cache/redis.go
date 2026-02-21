package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis - Redis 클라이언트 초기화
func InitRedis(addr string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}

// InitRedisWithPool - 연결 풀 설정으로 초기화
func InitRedisWithPool(addr, password string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		PoolTimeout:  4 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}

// GetRedisClient - 클라이언트 반환
func GetRedisClient() *redis.Client {
	return redisClient
}

// CloseRedis - 연결 종료
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// SetRedis - Redis에 저장
func SetRedis(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, data, expiration).Err()
}

// GetRedis - Redis에서 조회
func GetRedis(ctx context.Context, key string, dest interface{}) error {
	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// DeleteRedis - Redis에서 삭제
func DeleteRedis(ctx context.Context, key string) error {
	return redisClient.Del(ctx, key).Err()
}

// ExistsRedis - 키 존재 여부
func ExistsRedis(ctx context.Context, key string) (bool, error) {
	n, err := redisClient.Exists(ctx, key).Result()
	return n > 0, err
}

// DeleteByPattern - 패턴으로 삭제
func DeleteByPattern(ctx context.Context, pattern string) error {
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		if err := redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

// SetHash - 해시로 저장
func SetHash(ctx context.Context, key string, fields map[string]interface{}) error {
	return redisClient.HSet(ctx, key, fields).Err()
}

// GetHash - 해시 조회
func GetHash(ctx context.Context, key string) (map[string]string, error) {
	return redisClient.HGetAll(ctx, key).Result()
}

// GetHashField - 특정 필드만 조회
func GetHashField(ctx context.Context, key, field string) (string, error) {
	return redisClient.HGet(ctx, key, field).Result()
}
