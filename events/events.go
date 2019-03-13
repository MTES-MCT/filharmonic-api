package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type EventsManager struct {
	ws *melody.Melody
}

func New() *EventsManager {
	manager := &EventsManager{
		ws: melody.New(),
	}
	manager.ws.HandleMessage(manager.handleMessage)
	return manager
}

type Event struct {
	Type    EventType `json:"type"`
	Payload Payload   `json:"payload"`
}

type EventType string

const (
	// Client events
	EventConnect     EventType = "connect"
	EventDisonnect   EventType = "disconnect"
	EventSubscribe   EventType = "subscribe"
	EventUnsubscribe EventType = "unsubscribe"

	// Server events
	EventResourceUpdated EventType = "resource_updated"
)

type Payload struct {
	// connect / disconnect
	UserId int64 `json:"user_id,omitempty"`

	// subscribe / unsubscribe
	Resource   string `json:"resource,omitempty"`
	ResourceId int64  `json:"resource_id,omitempty"`
}

func (em *EventsManager) handleMessage(s *melody.Session, msg []byte) {
	event := Event{}
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Error().Err(err).Bytes("msg", msg).Msg("could not unmarshal msg")
		return
	}
	log.Debug().Interface("event", event).Msg("new ws event")
	switch event.Type {
	case EventConnect:
		s.Set("user_id", event.Payload.UserId)
	case EventDisonnect:
		s.Set("user_id", int64(0))
	case EventSubscribe:
		key := event.Payload.Resource
		if event.Payload.ResourceId > 0 {
			key += ":" + strconv.FormatInt(event.Payload.ResourceId, 10)
		}
		s.Set(key, true)
	case EventUnsubscribe:
		key := event.Payload.Resource
		if event.Payload.ResourceId > 0 {
			key += ":" + strconv.FormatInt(event.Payload.ResourceId, 10)
		}
		s.Set(key, false)
	default:
		log.Error().Interface("event", event).Msg("unknown event type")
	}
}

func (em *EventsManager) DispatchUpdatedResource(ctx *domain.UserContext, resource string, id int64) error {
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
		user_id, ok := s.Get("user_id")
		if !ok || user_id.(int64) == 0 {
			return false
		}
		value, ok := s.Get(resourceKey)
		return ok && value.(bool) && user_id != ctx.User.Id
	})
	return err
}

func (em *EventsManager) DispatchUpdatedNotifications(ids []int64) error {
	fmt.Printf("dispatch notif ids= %+v\n", ids)
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
		user_id, ok := s.Get("user_id")
		if !ok || user_id.(int64) == 0 {
			return false
		}
		for _, id := range ids {
			if id == user_id.(int64) {
				return true
			}
		}
		return false
	})
	return err
}

func (em *EventsManager) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return em.ws.HandleRequest(w, r)
}
