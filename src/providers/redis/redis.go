package redis

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/scoville/scvl/src/engine"
)

// NewClient creates and returns redis client
func NewClient() (r engine.RedisClient, err error) {
	c, err := redigo.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	r = &redisClient{c}
	return
}

type redisClient struct {
	c redigo.Conn
}

const ttl = 180

func (r *redisClient) SetURL(slug, url string) {
	r.c.Send("MULTI")
	r.c.Send("SET", r.urlKey(slug), url)
	r.c.Send("EXPIRE", r.urlKey(slug), ttl)
	r.c.Do("EXEC")
}

func (r *redisClient) GetURL(slug string) string {
	url, err := redigo.String(r.c.Do("GET", r.urlKey(slug)))
	if err != nil {
		return ""
	}
	return url
}

func (r *redisClient) DeleteURL(slug string) {
	r.c.Send("DEL", r.urlKey(slug))
}

func (r *redisClient) SetOGPID(slug string, id int) {
	r.c.Send("MULTI")
	r.c.Send("SET", r.ogpIDKey(slug), fmt.Sprintf("%d", id))
	r.c.Send("EXPIRE", r.urlKey(slug), ttl)
	r.c.Do("EXEC")
}

func (r *redisClient) GetOGPID(slug string) int {
	ogpID, err := redigo.Int(r.c.Do("GET", r.ogpIDKey(slug)))
	if err != nil {
		return 0
	}
	return ogpID
}

func (r *redisClient) DeleteOGPID(slug string) {
	r.c.Send("DEL", r.ogpIDKey(slug))
}

func (r *redisClient) Close() (err error) {
	return r.c.Close()
}

func (r *redisClient) SetDatabase(db uint) {
	r.c.Send("SELECT", db)
}

func (r *redisClient) urlKey(slug string) string {
	return "url_" + slug
}

func (r *redisClient) ogpIDKey(slug string) string {
	return "ogp_id_" + slug
}
