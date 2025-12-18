package types

import (
	"time"
)

type Model struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	CID       string            `json:"cid"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string  `json:"metadata"`
}

