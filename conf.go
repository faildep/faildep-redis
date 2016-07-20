package redis

import "time"

type RedisConfig struct {
	Basis     RedisBasic
	Resilient RedisResilient
}

type RedisBasic struct {
	MaxIdle      int
	MaxActive    int
	IdleTimeout  time.Duration
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type RedisResilient struct {
	ActiveReqThreshold      uint64
	ActiveReqCountWindow    time.Duration
	SuccessiveFailThreshold uint
	TrippedBaseTime         time.Duration
	TrippedTimeoutMax       time.Duration
	MaxServerPick           uint
	RetryMaxInterval        time.Duration
}
