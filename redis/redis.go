package redis

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Addr     string `json:"addr" default:"localhost:6379"`
	Password string `json:"password" default:"filharmonic"`
	Database int    `json:"database" default:"0"`
}

type Client = redis.Client

func New(config Config) (*Client, error) {
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
	return redisClient, nil
}
