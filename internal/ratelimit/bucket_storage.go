package ratelimit

import (
	"sync"
	"time"
)

type BucketStorage interface {
	GetBucket(key string) (*Bucket, bool)
	AddBucket(key string, b *Bucket)
	DeleteBucket(key string)
}

type MemStorage struct {
	bucketMap sync.Map
}

func (s *MemStorage) GetBucket(key string) (*Bucket, bool) {
	b, ok := s.bucketMap.Load(key)
	if !ok {
		return nil, false
	}

	return b.(*Bucket), true
}

func (s *MemStorage) AddBucket(key string, b *Bucket) {
	s.bucketMap.Store(key, b)
}

func (s *MemStorage) DeleteBucket(key string) {
	s.bucketMap.Delete(key)
}

// Clean execute in goroutine with some interval.
func (s *MemStorage) Clean(now time.Time, period time.Duration) {
	s.bucketMap.Range(func(k any, v any) bool {
		if v != nil {
			b := v.(*Bucket)
			if b.IsExpired(now, period) {
				s.bucketMap.Delete(k)
			}
		}

		return true
	})
}
