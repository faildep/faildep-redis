package go_resilient_redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/context"
)

var (
	ErrSet     error = errors.New("Set Redis failure")
	ErrKeyMiss       = errors.New("Set Key missed")
)

func (r *ResilientRedis) Get(ctx context.Context, key string) (interface{}, error) {
	return r.Do("GET", key)
}

func (r *ResilientRedis) Set(ctx context.Context, key string, value interface{}, expire int32) error {
	var (
		reply interface{}
		err   error
	)
	if expire > 0 {
		reply, err = r.Do("SETEX", key, expire, value)
	} else {
		reply, err = r.Do("SET", key, value)
	}
	if err != nil {
		return err
	}
	if reply != "OK" {
		return ErrSet
	}
	return nil
}

func (r *ResilientRedis) Delete(ctx context.Context, key string) error {
	reply, err := r.Do("DEL", key)
	num, err := redis.Int(reply, err)
	if err != nil {
		return err
	} else if num <= 0 {
		return ErrKeyMiss
	}
	return nil
}
