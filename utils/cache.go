package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var (
	redis_cache *redis.Client
	ctx         context.Context
)

func InitCache() error {
	ctx = context.Background()
	address := os.Getenv("REDIS_ADDRESS")
	if address == "" {
		address = "redis" // Giá trị mặc định nếu không có biến môi trường HOST
	}
	redis_cache = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if pong, err := redis_cache.Ping(ctx).Result(); err != nil {
		fmt.Println("Error ---> ", err)
		return errors.New("Redis: Kết nối không thành công!! " + err.Error())
	} else {
		fmt.Println("Redis: Kết nối thành công!! ", pong)
		return nil
	}
}
func SetCache(key, value string, duration int) error {
	fmt.Println("Cache ---> key --> value --> duration", key, value, time.Minute*time.Duration(duration))
	return redis_cache.Set(ctx, key, value, time.Minute*time.Duration(duration)).Err()
}
func SetCacheInterface(key string, data interface{}, duration int) error {
	value, err := json.Marshal(data)
	if err != nil {
		fmt.Println("SetCacheInterface: Fail to cache key")
		return err
	}
	return redis_cache.Set(ctx, key, value, time.Minute*time.Duration(duration)).Err()
}
func GetCache(key string, out interface{}) error {
	if value, err := redis_cache.Get(ctx, key).Result(); err == nil {
		return json.Unmarshal([]byte(value), &out)
	} else {
		return err
	}
}

func DelCache(key string) error {
	_, err := redis_cache.Del(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else {
		return err
	}
}
func CheckTTL(key string) (time.Duration, error) {
	ttl, err := redis_cache.TTL(ctx, key).Result()
	if err == redis.Nil {
		return ttl, nil
	} else {
		return ttl, err
	}
}

//func DelAllCache() error  {
//	value, err := redis_cache
//}
