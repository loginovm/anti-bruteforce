package ratelimit

import "time"

type Calc struct {
	period  time.Duration
	time    currentTime
	storage BucketStorage
}

func NewCalc(period time.Duration, storage BucketStorage) *Calc {
	return &Calc{
		period:  period,
		time:    &curTimeReal{},
		storage: storage,
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
	if b.IsExpired(now, c.period) {
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
