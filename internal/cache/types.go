package cache

import "time"

type Memory struct {
	Query     string
	Response  string
	CreatedAt time.Time
}
