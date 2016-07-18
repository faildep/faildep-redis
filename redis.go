package go_resilient_redis

import (
	"github.com/lysu/slb"
	"time"
)

type ResilientRedis struct {
	pool *redisPool
	lb   *slb.LoadBalancer
}

func NewResilientRedis(servers []string, conf RedisConfig) (*ResilientRedis, error) {
	p, err := newRedisPool(servers, conf.Basis)
	if err != nil {
		return nil, err
	}
	l := slb.NewLoadBalancer(servers,
		slb.WithBulkhead(conf.Resilient.ActiveReqThreshold, conf.Resilient.ActiveReqCountWindow),
		slb.WithCircuitBreaker(
			conf.Resilient.SuccessiveFailThreshold, conf.Resilient.TrippedBaseTime,
			conf.Resilient.TrippedTimeoutMax, slb.DecorrelatedJittered,
		),
		slb.WithRetry(
			conf.Resilient.MaxServerPick, 1, 0*time.Millisecond,
			conf.Resilient.RetryMaxInterval, slb.DecorrelatedJittered,
		),
	)
	return &ResilientRedis{pool: p, lb: l}, nil
}

func (r *ResilientRedis) Do(commandName string, args ...interface{}) (interface{}, error) {
	var reply interface{}
	err := r.lb.Submit(func(node *slb.Node) error {
		instance, err := r.pool.get(node.Server)
		if err != nil {
			return err
		}
		conn := instance.Get()
		defer func() {
			if conn != nil {
				conn.Close()
			}
		}()
		reply, err = conn.Do(commandName, args...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return reply, nil
}
