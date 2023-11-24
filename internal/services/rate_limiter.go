package services

import (
	"time"
	"webhook/internal/config"

	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	Limiter ratelimit.Limiter
}

func NewRateLimiter(config config.Configuration) *RateLimiter {
	return &RateLimiter{
		Limiter: ratelimit.New(config.Requests.PerSecond),
	}
}

func (r *RateLimiter) Take() time.Time {
	return r.Limiter.Take()
}
