package ratelimit

import "time"

type currentTime interface {
	Now() time.Time
}

type curTimeReal struct{}

func (t *curTimeReal) Now() time.Time {
	return time.Now()
}

type curTimeFake struct {
	now time.Time
}

func (t *curTimeFake) Now() time.Time {
	return t.now
}
