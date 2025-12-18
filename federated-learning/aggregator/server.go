package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/atlas/federated-learning/protocols"
)

type Aggregator struct {
	clients      map[string]*ClientState
	protocol     *protocols.FLProtocol
	jobID        string
	round        int
	nodeID       string // This node's ID (for distributed aggregation)
	isAggregator bool   // Whether this node is acting as aggregator
}

type ClientState struct {
	NodeID      string
	Gradients   []float64
	Ready       bool
	Timestamp   time.Time
	Contribution float64 // Contribution weight for fairness
}

func NewAggregator(jobID string, nodeID string, ipfsAPIURL string) *Aggregator {
	return &Aggregator{
		clients:      make(map[string]*ClientState),
		protocol:     protocols.NewFLProtocol(ipfsAPIURL, "aggregator"),
		jobID:        jobID,
		nodeID:       nodeID,
		isAggregator: false,
		round:        0,
	}
}

func (a *Aggregator) BecomeAggregator() error {
	a.isAggregator = true
	topic := fmt.Sprintf("/atlas/fl/aggregator/%s", a.jobID)
	
	announcement := map[string]interface{}{
		"node_id":    a.nodeID,
		"job_id":     a.jobID,
		"timestamp":  time.Now().Unix(),
		"action":     "become_aggregator",
	}
	
	data, err := json.Marshal(announcement)
	if err != nil {
		return fmt.Errorf("failed to marshal announcement: %w", err)
	}
	
	if err := a.protocol.GetAPI().PubSubPublish(topic, string(data)); err != nil {
		return fmt.Errorf("failed to publish aggregator announcement: %w", err)
	}
	
	return nil
}

func (a *Aggregator) Start(ctx context.Context) error {
	topic := fmt.Sprintf("/atlas/fl/gradients/%s", a.jobID)
	
	sub, err := a.protocol.GetAPI().PubSubSubscribe(topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to gradient topic: %w", err)
	}
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-sub:
				var gradientMsg protocols.GradientMessage
				if err := json.Unmarshal(msg.Data, &gradientMsg); err != nil {
					continue
				}
				
				if gradientMsg.Round == a.round {
					if err := a.ReceiveGradients(gradientMsg.NodeID, gradientMsg.Gradients); err != nil {
						continue
					}
				}
			}
		}
	}()
	
	return nil
}

func (a *Aggregator) ReceiveGradients(nodeID string, gradients []float64) error {
	if state, ok := a.clients[nodeID]; ok {
		state.Gradients = gradients
		state.Ready = true
		state.Timestamp = time.Now()
	} else {
		a.clients[nodeID] = &ClientState{
			NodeID:    nodeID,
			Gradients: gradients,
			Ready:     true,
			Timestamp: time.Now(),
		}
	}
	return nil
}

func (a *Aggregator) Aggregate() ([]float64, error) {
	var gradientsList [][]float64
	var weights []float64
	
	for _, state := range a.clients {
		if state.Ready {
			gradientsList = append(gradientsList, state.Gradients)
			weights = append(weights, 1.0)
		}
	}
	
	if len(gradientsList) == 0 {
		return nil, fmt.Errorf("no gradients to aggregate")
	}
	
	aggregated, err := FederatedAveraging(gradientsList, weights)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	
	for _, state := range a.clients {
		state.Ready = false
	}
	
	a.round++
	return aggregated, nil
}

