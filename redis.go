package redis

import (
	"github.com/faildep/faildep"
	"time"
)

type FailDepRedis struct {
	pool    *redisPool
	failDep *faildep.FailDep
}

func NewFailDepRedis(servers []string, conf RedisConfig) (*FailDepRedis, error) {
	p, err := newRedisPool(servers, conf.Basis)
	if err != nil {
		return nil, err
	}
	f := faildep.NewFailDep(servers,
		faildep.WithBulkhead(conf.Resilient.ActiveReqThreshold, conf.Resilient.ActiveReqCountWindow),
		faildep.WithCircuitBreaker(
			conf.Resilient.SuccessiveFailThreshold, conf.Resilient.TrippedBaseTime,
			conf.Resilient.TrippedTimeoutMax, faildep.DecorrelatedJittered,
		),
		faildep.WithRetry(
			conf.Resilient.MaxServerPick, 1, 0*time.Millisecond,
			conf.Resilient.RetryMaxInterval, faildep.DecorrelatedJittered,
		),
	)
	return &FailDepRedis{pool: p, failDep: f}, nil
}

func (r *FailDepRedis) Do(commandName string, args ...interface{}) (interface{}, error) {
	var reply interface{}
	err := r.failDep.Do(func(res *faildep.Resource) error {
		instance, err := r.pool.get(res.Server)
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
