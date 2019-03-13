package sessions

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/redis"
	goredis "github.com/go-redis/redis"
	"github.com/gofrs/uuid"
)

const Prefix = "filharmonic:sessions:"

type Config struct {
	Expiration time.Duration `json:"expiration" default:"604800s"` // 7 days
}

type RedisSessions struct {
	config Config
	client *redis.Client
}

func NewRedis(config Config, redisClient *redis.Client) *RedisSessions {
	return &RedisSessions{
		config: config,
		client: redisClient,
	}
}

func (redisSessions *RedisSessions) Get(sessionToken string) (int64, error) {
	id, err := redisSessions.client.Get(Prefix + sessionToken).Int64()
	if err != nil {
		if err == goredis.Nil {
			return int64(0), nil
		}
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
