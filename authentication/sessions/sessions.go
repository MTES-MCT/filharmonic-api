package sessions

type Sessions interface {
	Get(sessionToken string) (int64, error)
	Add(userId int64) (string, error)
	Set(sessionToken string, userId int64) error
	Delete(sessionToken string) error
	Close() error
}
