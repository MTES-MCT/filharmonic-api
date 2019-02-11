package sessions

import (
	"github.com/gofrs/uuid"
)

type MemorySessions struct {
	sessions map[string]int64
}

func New() *MemorySessions {
	return &MemorySessions{
		sessions: make(map[string]int64),
	}
}

func (memorySessions *MemorySessions) Get(sessionToken string) int64 {
	return memorySessions.sessions[sessionToken]
}

func (memorySessions *MemorySessions) Add(userId int64) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken := id.String()
	memorySessions.sessions[sessionToken] = userId
	return sessionToken, nil
}

func (memorySessions *MemorySessions) Set(sessionToken string, userId int64) {
	memorySessions.sessions[sessionToken] = userId
}

func (memorySessions *MemorySessions) Delete(sessionToken string) {
	delete(memorySessions.sessions, sessionToken)
}
