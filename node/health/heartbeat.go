package health

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-ipfs-api"
)

type Monitor struct {
	nodeID    string
	ipfsAPI   *api.Shell
	lastBeat  time.Time
}

func NewMonitor() *Monitor {
	return &Monitor{
		nodeID:   generateNodeID(),
		ipfsAPI:  api.NewShell("localhost:5001"),
		lastBeat: time.Now(),
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			m.sendHeartbeat()
		}
	}
}

func (m *Monitor) sendHeartbeat() {
	topic := fmt.Sprintf("/atlas/heartbeat/%s", m.nodeID)
	message := fmt.Sprintf(`{"node_id":"%s","timestamp":"%s"}`, m.nodeID, time.Now().Format(time.RFC3339))
	m.ipfsAPI.PubSubPublish(topic, message)
	m.lastBeat = time.Now()
}

func generateNodeID() string {
	return fmt.Sprintf("node-%d", time.Now().UnixNano())
}

