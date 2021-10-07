package storage

import "time"

type Config struct {
	Bucket          string
	PresignDuration time.Duration
}
