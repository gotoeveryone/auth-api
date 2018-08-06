package cache

// Client is operationg of redis connection.
type Client interface {
	Connect() error
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (interface{}, error)
	Delete(key string) (interface{}, error)
	Expire(key string, expire int) (interface{}, error)
	SetWithExpire(key string, expire int, value interface{}) error
}
