package storage

import (
	"time"
)

type Disk interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	SetWithMaxAge(string, []byte, time.Duration) error
	Remove(string) error
}
