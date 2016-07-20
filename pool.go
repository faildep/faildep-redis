package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type redisPool struct {
	clients map[string]*redis.Pool
}

func newRedisPool(servers []string, conf RedisBasic) (*redisPool, error) {
	cs := make(map[string]*redis.Pool, len(servers))
	for _, srv := range servers {
		p := &redis.Pool{
			MaxIdle:     conf.MaxIdle,
			MaxActive:   conf.MaxActive,
			IdleTimeout: conf.IdleTimeout,
			Dial: func(serverAddr string) func() (redis.Conn, error) {
				return func() (redis.Conn, error) {
					c, err := redis.DialTimeout("tcp", serverAddr,
						conf.ConnTimeout,
						conf.ReadTimeout,
						conf.WriteTimeout)
					if err != nil {
						return nil, err
					}
					return c, err
				}
			}(srv),
		}
		cs[srv] = p
	}
	return &redisPool{clients: cs}, nil
}

func (p *redisPool) get(server string) (*redis.Pool, error) {
	c, ok := p.clients[server]
	if !ok {
		return nil, fmt.Errorf("Mc Server %s not found", server)
	}
	return c, nil
}
