package events

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/redis"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type Config struct {
	TTL    time.Duration `default:"168h"` // 1 week
	Prefix string        `default:"filharmonic:ws:"`
}

type RedisEventsManager struct {
	config      Config
	ws          *melody.Melody
	redisClient *redis.Client
}

func New(config Config, redisClient *redis.Client) *RedisEventsManager {
	manager := &RedisEventsManager{
		ws:          melody.New(),
		redisClient: redisClient,
	}
	manager.ws.HandleMessage(manager.handleMessage)
	return manager
}

func (em *RedisEventsManager) handleMessage(s *melody.Session, msg []byte) {
	event := Event{}
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Error().Err(err).Bytes("msg", msg).Msg("could not unmarshal msg")
		return
	}
	log.Debug().Interface("event", event).Msg("new ws event")

	switch event.Type {
	case EventConnect:
		em.redisClient.Set(em.sessionPrefix(s)+"user_id", event.Payload.UserId, em.config.TTL)
	case EventDisonnect:
		em.redisClient.Del(em.sessionPrefix(s) + "user_id")
	case EventSubscribe:
		key := event.Payload.Resource
		if event.Payload.ResourceId > 0 {
			key += ":" + strconv.FormatInt(event.Payload.ResourceId, 10)
		}
		em.redisClient.Set(em.sessionPrefix(s)+key, true, em.config.TTL)
	case EventUnsubscribe:
		key := event.Payload.Resource
		if event.Payload.ResourceId > 0 {
			key += ":" + strconv.FormatInt(event.Payload.ResourceId, 10)
		}
		em.redisClient.Del(em.sessionPrefix(s) + key)
	default:
		log.Error().Interface("event", event).Msg("unknown event type")
	}
}

func (em *RedisEventsManager) DispatchUpdatedResource(ctx *domain.UserContext, resource string, id int64) error {
	resourceKey := resource + ":" + strconv.FormatInt(id, 10)
	event := Event{
		Type: EventResourceUpdated,
		Payload: Payload{
			Resource:   resource,
			ResourceId: id,
		},
	}
	eventStr, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = em.ws.BroadcastFilter(eventStr, func(s *melody.Session) bool {
		userIdStr, err := em.redisClient.Get(em.sessionPrefix(s) + "user_id").Result()
		if err != nil {
			log.Error().Err(err).Msg("could not get user_id key")
			return false
		}
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("could not parse user_id")
			return false
		}
		value, err := em.redisClient.Exists(em.sessionPrefix(s) + resourceKey).Result()
		if err != nil {
			log.Error().Err(err).Msg("could not get resource key")
			return false
		}
		return value == 1 && userId != ctx.User.Id
	})
	return err
}

func (em *RedisEventsManager) DispatchUpdatedNotifications(ids []int64) error {
	event := Event{
		Type: EventResourceUpdated,
		Payload: Payload{
			Resource: "notifications",
		},
	}
	eventStr, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = em.ws.BroadcastFilter(eventStr, func(s *melody.Session) bool {

		userIdStr, err := em.redisClient.Get(em.sessionPrefix(s) + "user_id").Result()
		if err != nil {
			log.Error().Err(err).Msg("could not get user_id key")
			return false
		}
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("could not parse user_id")
			return false
		}
		for _, id := range ids {
			if id == userId {
				return true
			}
		}
		return false
	})
	return err
}

func (em *RedisEventsManager) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return em.ws.HandleRequest(w, r)
}

func (em *RedisEventsManager) sessionPrefix(s *melody.Session) string {
	return em.config.Prefix + s.Request.Header.Get("Sec-Websocket-Key") + ":"
}
