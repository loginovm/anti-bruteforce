package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalc(t *testing.T) {
	t.Run("all cases workflow", func(t *testing.T) {
		fakeTime := &curTimeFake{now: time.Now()}
		sut := Calc{
			period:  time.Minute,
			time:    fakeTime,
			storage: &MemStorage{},
		}
		key := "k"
		limit := 2
		actual := sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)

		// limit exceeded
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, false, actual)

		// minute passed
		fakeTime.now = fakeTime.now.Add(time.Minute + time.Second)
		_ = sut.TryIncrement(key, limit)
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)

		// limit exceeded
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, false, actual)

		// 59 sec passed
		fakeTime.now = fakeTime.now.Add(59 * time.Second)
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, false, actual)

		// 2 sec passed
		fakeTime.now = fakeTime.now.Add(2 * time.Second)
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)
	})

	t.Run("bucket reset", func(t *testing.T) {
		fakeTime := &curTimeFake{now: time.Now()}
		sut := Calc{
			period:  time.Minute,
			time:    fakeTime,
			storage: &MemStorage{},
		}
		key := "k"
		limit := 2
		_ = sut.TryIncrement(key, limit)
		actual := sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)

		// limit exceeded
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, false, actual)

		sut.storage.DeleteBucket(key)
		_ = sut.TryIncrement(key, limit)
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, true, actual)

		// limit exceeded
		actual = sut.TryIncrement(key, limit)
		require.Equal(t, false, actual)
	})
}
