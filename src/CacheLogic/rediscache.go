// Package CacheLogic holds the interface for the cache handlers and all the implementations
package CacheLogic

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

// RedisCache encapsulates the basic redis get , insert and key exists operations , accepts a redis client pointer
type RedisCache struct {
	Rdb *redis.Client
}

var ctx = context.Background()

//Get returns the string value of the given key
//throws exception if the key is not found
func (m RedisCache) Get(key string) string {
	val, err := m.Rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

//Insert sets the key's value in the redis cache with val. ttl is "time to live" in minutes
func (m RedisCache) Insert(key string, val *string, ttl int) {
	err := m.Rdb.Set(ctx, key, *val, time.Duration(ttl)*time.Minute).Err()
	if err != nil {
		panic(err)
	}
}

//KeyExists returns true if the key has a value in redis, false if not found.
func (m RedisCache) KeyExists(key string) bool {
	_, err := m.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	} else {
		return true
	}
}
