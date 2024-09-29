package ratelimit

import "time"

type Calc struct {
	resetCounterAfter time.Duration
	time              currentTime
	storage           BucketStorage
}

func NewCalc(resetCounterAfter time.Duration, storage BucketStorage) *Calc {
	return &Calc{
		resetCounterAfter: resetCounterAfter,
		time:              &curTimeReal{},
		storage:           storage,
	}
}

// TryIncrement increments attempts counter for 'key' and returns false if 'limit' has been exceeded.
func (c *Calc) TryIncrement(key string, limit int) bool {
	b, ok := c.storage.GetBucket(key)
	if !ok {
		// this is first time key accessed
		b = &Bucket{
			startDate: c.time.Now(),
			counter:   1,
		}
		c.storage.AddBucket(key, b)
		return true
	}
	now := c.time.Now()
	if b.IsExpired(now, c.resetCounterAfter) {
		// control period elapsed so reset counter
		b.startDate = now
		b.counter = 1
		return true
	}

	// we are within control period so increment and check limit
	if b.counter < limit {
		b.counter++
		return true
	}

	return false
}

func (c *Calc) ResetBucket(key string) {
	c.storage.DeleteBucket(key)
}
