package session

import "time"

const (
	DefaultIdleTimeout = 30 * time.Minute
	DefaultHardTimeout = 24 * time.Hour
)

func IsExpired(hard_deadline, idle_deadline time.Time) bool {
	now := time.Now()
	return hard_deadline.Before(now) || idle_deadline.Before(now)
}
