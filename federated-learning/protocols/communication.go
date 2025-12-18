package protocols

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/ipfs/go-ipfs-api"
)

type FLProtocol struct {
	api    *api.Shell
	nodeID string
}

func NewFLProtocol(apiURL string, nodeID string) *FLProtocol {
	return &FLProtocol{
		api:    api.NewShell(apiURL),
		nodeID: nodeID,
	}
}

func (p *FLProtocol) GetAPI() *api.Shell {
	return p.api
}

type GradientMessage struct {
	NodeID    string    `json:"node_id"`
	JobID     string    `json:"job_id"`
	Round     int       `json:"round"`
	Gradients []float64 `json:"gradients"`
	Timestamp string    `json:"timestamp"`
}

func (p *FLProtocol) SendGradients(ctx context.Context, jobID string, round int, gradients []float64) error {
	topic := fmt.Sprintf("/atlas/fl/gradients/%s", jobID)
	
	msg := GradientMessage{
		NodeID:    p.nodeID,
		JobID:     jobID,
		Round:     round,
		Gradients: gradients,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.api.PubSubPublish(topic, string(data))
}

func (p *FLProtocol) ReceiveModel(ctx context.Context, jobID string, handler func([]float64)) error {
	topic := fmt.Sprintf("/atlas/fl/model/%s", jobID)
	
	sub, err := p.api.PubSubSubscribe(topic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-sub:
				var modelUpdate struct {
					Weights []float64 `json:"weights"`
				}
				if err := json.Unmarshal(msg.Data, &modelUpdate); err == nil {
					handler(modelUpdate.Weights)
				}
			}
		}
	}()

	return nil
}

