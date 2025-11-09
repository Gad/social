package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration, error)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
