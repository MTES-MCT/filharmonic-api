package events

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
)

type EventsManager interface {
	DispatchUpdatedResource(ctx *domain.UserContext, resource string, id int64) error
	DispatchUpdatedNotifications(ids []int64) error
	HandleRequest(w http.ResponseWriter, r *http.Request) error
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

	// subscribe / unsubscribe / resource_updated
	Resource   string `json:"resource,omitempty"`
	ResourceId int64  `json:"resource_id,omitempty"`
}
