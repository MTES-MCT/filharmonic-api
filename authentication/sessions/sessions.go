package sessions

import "strconv"

var sessions = make(map[string]int64)

func Get(sessionToken string) int64 {
	return sessions[sessionToken]
}

func New(userId int64) string {
	// TODO random token
	sessionToken := "token-" + strconv.FormatInt(userId, 10)
	sessions[sessionToken] = userId
	return sessionToken
}

func Set(sessionToken string, userId int64) {
	sessions[sessionToken] = userId
}

func Delete(sessionToken string) {
	delete(sessions, sessionToken)
}
