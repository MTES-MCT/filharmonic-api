package sessions

import (
	"github.com/gofrs/uuid"
)

type MemorySessions struct {
	sessions map[string]int64
}

func NewMemory() *MemorySessions {
	return &MemorySessions{
		sessions: make(map[string]int64),
	}
}

func (memorySessions *MemorySessions) Get(sessionToken string) (int64, error) {
	return memorySessions.sessions[sessionToken], nil
}

func (memorySessions *MemorySessions) Add(userId int64) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken := id.String()
	err = memorySessions.Set(sessionToken, userId)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (memorySessions *MemorySessions) Set(sessionToken string, userId int64) error {
	memorySessions.sessions[sessionToken] = userId
	return nil
}

func (memorySessions *MemorySessions) Delete(sessionToken string) error {
	delete(memorySessions.sessions, sessionToken)
	return nil
}
