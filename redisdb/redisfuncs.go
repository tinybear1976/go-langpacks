package redisdb

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	serverTags map[string]*redis.Pool = make(map[string]*redis.Pool)
)

func New(serverTag, server, password string, dbnum int) {
	redisPool := &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 240 * time.Second,
		MaxActive:   1000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialDatabase(dbnum))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	serverTags[serverTag] = redisPool
}

func Destroy() {
	for k := range serverTags {
		delete(serverTags, k)
	}
}

func getConn(serverTag string) (*redis.Conn, error) {
	pool, ok := serverTags[serverTag]
	if !ok {
		return nil, errors.New(fmt.Sprintf("redis[%s] not existing", serverTag))
	}
	redisClient := pool.Get()
	return &redisClient, nil
}

func GET(serverTag, key string) (string, error) {
	rc, err := getConn(serverTag)
	if err != nil {
		return "", err
	}
	defer (*rc).Close()
	val, err := redis.String((*rc).Do("GET", key))
	return val, err
}

func SET(serverTag, key, value string) error {
	rc, err := getConn(serverTag)
	if err != nil {
		return err
	}
	defer (*rc).Close()
	_, err = (*rc).Do("SET", key, value)
	return err
}

func DEL(serverTag string, keys ...interface{}) error {
	rc, err := getConn(serverTag)
	if err != nil {
		return err
	}
	defer (*rc).Close()
	_, err = (*rc).Do("DEL", keys...)
	return err
}

func KEYS(serverTag, query string) ([]string, error) {
	rc, err := getConn(serverTag)
	if err != nil {
		return nil, err
	}
	defer (*rc).Close()
	keys, err := redis.Strings((*rc).Do("KEYS", query))
	return keys, err
}

func HMSET(serverTag string, params ...interface{}) error {
	rc, err := getConn(serverTag)
	if err != nil {
		return err
	}
	defer (*rc).Close()
	_, err = (*rc).Do("HMSET", params...)
	return err
}

func HMGET(serverTag string, params ...interface{}) ([]string, error) {
	rc, err := getConn(serverTag)
	if err != nil {
		return nil, err
	}
	defer (*rc).Close()
	vals, err := redis.Strings((*rc).Do("HMGET", params...))
	return vals, err
}

func HGETALL(serverTag, key string) (map[string]string, error) {
	rc, err := getConn(serverTag)
	if err != nil {
		return nil, err
	}
	defer (*rc).Close()
	ret, err := redis.StringMap((*rc).Do("HGETALL", key))
	return ret, err
}

func HDEL(serverTag string, params ...interface{}) error {
	rc, err := getConn(serverTag)
	if err != nil {
		return err
	}
	defer (*rc).Close()
	_, err = (*rc).Do("HDEL", params...)
	return err
}

func EXISTS(serverTag, key string) (bool, error) {
	rc, err := getConn(serverTag)
	if err != nil {
		return false, err
	}
	defer (*rc).Close()
	res, err := redis.Int((*rc).Do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return res != 0, nil
}
