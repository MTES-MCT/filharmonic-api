package events

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/olahol/melody"
)

type StubEventsManager struct {
	ws *melody.Melody
}

func NewStub() *StubEventsManager {
	return &StubEventsManager{
		ws: melody.New(),
	}
}

func (em *StubEventsManager) DispatchUpdatedResource(ctx *domain.UserContext, resource string, id int64) error {
	return nil
}

func (em *StubEventsManager) DispatchUpdatedNotifications(ids []int64) error {
	return nil
}

func (em *StubEventsManager) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return em.ws.HandleRequest(w, r)
}
