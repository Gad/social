package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration, error) {
	rl.Lock()

	count, exists := rl.clients[ip]
	rl.Unlock()
	if !exists || count < rl.limit {
		rl.Lock()
		if !exists {
			go rl.resetCountAfterWindow(ip)
		}

		rl.clients[ip]++
		rl.Unlock()

		return true, 0, nil
	}
	return false, rl.window, nil
}

func (rl *FixedWindowRateLimiter) resetCountAfterWindow(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
