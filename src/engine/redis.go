package engine

// RedisClient is the interface for redis client
type RedisClient interface {
	SetURL(string, string)
	GetURL(string) string
	DeleteURL(string)
	SetOGPID(string, int)
	GetOGPID(string) int
	DeleteOGPID(string)
	Close() error
	SetDatabase(uint)
}
