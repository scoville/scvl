package main

import (
	redis "github.com/garyburd/redigo/redis"
)

func newRedisClient() (r *redisClient, err error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	r = &redisClient{c}
	return
}

type redisClient struct {
	c redis.Conn
}

func (r *redisClient) SetURL(slug, url string) {
	r.c.Do("SET", r.urlKey(slug), url)
}

func (r *redisClient) GetURL(slug string) string {
	url, err := redis.String(r.c.Do("GET", r.urlKey(slug)))
	if err != nil {
		return ""
	}
	return url
}

func (r *redisClient) DeleteURL(slug string) {
	r.c.Send("DEL", r.urlKey(slug))
}

func (r *redisClient) Close() (err error) {
	return r.c.Close()
}

func (r *redisClient) setDatabase(db uint) {
	r.c.Send("SELECT", db)
}

func (r *redisClient) urlKey(slug string) string {
	return "url_" + slug
}
