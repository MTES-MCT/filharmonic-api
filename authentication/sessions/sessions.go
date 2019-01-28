package sessions

type Sessions interface {
	Get(sessionToken string) int64
	Add(userId int64) string
	Set(sessionToken string, userId int64)
	Delete(sessionToken string)
}
