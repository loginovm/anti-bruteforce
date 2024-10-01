package ratelimit

import "time"

type Bucket struct {
	startDate time.Time
	counter   int
}

func (b *Bucket) IsExpired(now time.Time, maxAge time.Duration) bool {
	elapsed := now.Sub(b.startDate)
	return elapsed > maxAge
}
