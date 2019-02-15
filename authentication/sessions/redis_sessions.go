package sessions

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

const Prefix = "filharmonic:sessions:"

type RedisConfig struct {
	Addr       string        `json:"addr" default:"localhost:6379"`
	Password   string        `json:"password" default:"filharmonic"`
	Database   int           `json:"database" default:"0"`
	Expiration time.Duration `json:"expiration" default:"604800s"` // 7 days
}

type RedisSessions struct {
	config RedisConfig
	client *redis.Client
}

func NewRedis(config RedisConfig) (*RedisSessions, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Database,
	})
	log.Info().Msgf("connecting to redis endpoint on %s", config.Addr)
	err := redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("connected to redis")
	return &RedisSessions{
		config: config,
		client: redisClient,
	}, nil
}

func (redisSessions *RedisSessions) Get(sessionToken string) (int64, error) {
	id, err := redisSessions.client.Get(Prefix + sessionToken).Int64()
	if err != nil {
		return int64(0), err
	}
	err = redisSessions.client.Expire(Prefix+sessionToken, redisSessions.config.Expiration).Err()
	return id, err
}

func (redisSessions *RedisSessions) Add(userId int64) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken := id.String()
	err = redisSessions.Set(sessionToken, userId)
	return sessionToken, err
}

func (redisSessions *RedisSessions) Set(sessionToken string, userId int64) error {
	return redisSessions.client.Set(Prefix+sessionToken, userId, redisSessions.config.Expiration).Err()
}

func (redisSessions *RedisSessions) Delete(sessionToken string) error {
	return redisSessions.client.Del(Prefix + sessionToken).Err()
}

func (redisSessions *RedisSessions) Close() error {
	return redisSessions.client.Close()
}
