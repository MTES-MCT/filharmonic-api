package sessions

import "strconv"

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

func (memorySessions *MemorySessions) Add(userId int64) string {
	// TODO random token
	sessionToken := "token-" + strconv.FormatInt(userId, 10)
	memorySessions.sessions[sessionToken] = userId
	return sessionToken
}

func (memorySessions *MemorySessions) Set(sessionToken string, userId int64) {
	memorySessions.sessions[sessionToken] = userId
}

func (memorySessions *MemorySessions) Delete(sessionToken string) {
	delete(memorySessions.sessions, sessionToken)
}
