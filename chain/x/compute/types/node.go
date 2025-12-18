package types

import (
	"fmt"
	"time"
)

type Node struct {
	ID              string            `json:"id"`
	Address         string            `json:"address"`
	Status          string            `json:"status"`
	Resources       map[string]string `json:"resources"`
	Reputation      float64           `json:"reputation"`
	UptimePercent   float64           `json:"uptime_percent"`
	LastHeartbeat   time.Time         `json:"last_heartbeat"`
	RegisteredAt    time.Time         `json:"registered_at"`
	ActiveTasks     []string          `json:"active_tasks"`
}

func (n Node) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}
	if n.Address == "" {
		return fmt.Errorf("node address cannot be empty")
	}
	return nil
}

