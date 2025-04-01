package impl

import (
	"context"
	"est-proxy/src/config"
	"est-proxy/src/models"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

var refreshPeriod time.Duration = 5 * time.Second

type RedisClientImpl struct {
	client        *redis.Client
	redisDB       string
	redisURL      string
	redisPassword string
}

func NewRedisClientImpl() *RedisClientImpl {
	redisClient := RedisClientImpl{
		redisDB:       config.REDIS_DB,
		redisURL:      config.REDIS_URL,
		redisPassword: config.REDIS_PASSWORD,
	}

	if err := redisClient.connect(); err != nil {
		log.Fatalf("Connect redis err: %v", err)
	}
	return &redisClient
}

func (r *RedisClientImpl) AddUser(ctx context.Context, userKey string, user *models.RegisteredUser) error {
	value := fmt.Sprintf(
		"%s |#&^| %s |#&^| %s",
		user.Username,
		user.Email,
		user.PasswordHash)
	return r.client.Set(ctx, userKey, value, config.REDIS_EXPIRATION_TIME).Err()
}

func (r *RedisClientImpl) GetUser(ctx context.Context, userKey string) (*models.RegisteredUser, error) {
	val, err := r.client.Get(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	user := models.RegisteredUser{}
	_, err = fmt.Sscanf(val, "%s |#&^| %s |#&^| %s",
		&user.Username,
		&user.Email,
		&user.PasswordHash)
	if err != nil {
		log.Printf("Failed to scan user data ( %s ) with userKey %s: %v", val, userKey, err)
		return nil, err
	}

	return &user, nil
}

func (r *RedisClientImpl) RemoveUser(ctx context.Context, userKey string) error {
	_, err := r.client.Del(ctx, userKey).Result()
	if err != nil {
		return err
	}
	return err
}

func (r *RedisClientImpl) Refresh() {
	ctx := context.Background()

	for {
		if err := r.client.Ping(ctx).Err(); err != nil {
			log.Println("Lost connection to Redis, attempting to reconnect...")
			failedAttempts := 0

			for {
				time.Sleep(refreshPeriod)
				if err := r.connect(); err == nil {
					log.Println("Successfully reconnected to Redis")
					break
				}

				if failedAttempts++; failedAttempts > 5 {
					log.Fatalf("Failed to reconnect to Redis after %d attempts", failedAttempts)
				}
			}
		}

		time.Sleep(refreshPeriod)
	}
}

func (r *RedisClientImpl) Close() {
	err := r.client.Close()
	if err != nil {
		log.Printf("Failed to close redis client: %v", err)
	}
}

func (r *RedisClientImpl) connect() error {
	db, err := strconv.Atoi(r.redisDB)
	if err != nil {
		return err
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.redisURL,
		Password: r.redisPassword,
		DB:       db,
	})

	return r.client.Ping(context.Background()).Err()
}
