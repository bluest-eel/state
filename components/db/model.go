package db

import "time"

// KV ...
type KV struct {
	Key     string
	Value   []byte
	Created *time.Time
	Updated *time.Time
}
