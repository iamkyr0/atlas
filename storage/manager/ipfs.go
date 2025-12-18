package manager

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/ipfs/go-ipfs-api"
)

type IPFSManager struct {
	api          *api.Shell
	fallbackAPIs []*api.Shell // Fallback IPFS nodes
	maxRetries   int
	retryDelay   time.Duration
}

func NewIPFSManager(apiURL string, fallbackURLs ...string) *IPFSManager {
	manager := &IPFSManager{
		api:        api.NewShell(apiURL),
		maxRetries: 3,
		retryDelay: 2 * time.Second,
	}
	
	for _, url := range fallbackURLs {
		manager.fallbackAPIs = append(manager.fallbackAPIs, api.NewShell(url))
	}
	
	return manager
}

func (m *IPFSManager) CheckDataAvailability(ctx context.Context, cid string) (bool, error) {
	if m.checkCID(ctx, m.api, cid) {
		return true, nil
	}
	
	for _, fallbackAPI := range m.fallbackAPIs {
		if m.checkCID(ctx, fallbackAPI, cid) {
			return true, nil
		}
	}
	
	return false, fmt.Errorf("data not available: %s", cid)
}

func (m *IPFSManager) checkCID(ctx context.Context, apiShell *api.Shell, cid string) bool {
	_, err := apiShell.ObjectStat(cid)
	return err == nil
}

func (m *IPFSManager) AddFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	cid, err := m.api.Add(file)
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (m *IPFSManager) GetFile(cid string, outputPath string) error {
	return m.GetFileWithFallback(cid, outputPath)
}

func (m *IPFSManager) GetFileWithFallback(cid string, outputPath string) error {
	var lastErr error
	
	err := m.api.Get(cid, outputPath)
	if err == nil {
		return nil
	}
	lastErr = err
	
	for i, fallbackAPI := range m.fallbackAPIs {
		delay := m.retryDelay * time.Duration(1<<uint(i))
		time.Sleep(delay)
		
		err := fallbackAPI.Get(cid, outputPath)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	
	return fmt.Errorf("failed to get file from all IPFS nodes: %w", lastErr)
}

func (m *IPFSManager) Pin(cid string) error {
	return m.api.Pin(cid)
}

func (m *IPFSManager) SetFallbackNodes(urls ...string) {
	m.fallbackAPIs = []*api.Shell{}
	for _, url := range urls {
		m.fallbackAPIs = append(m.fallbackAPIs, api.NewShell(url))
	}
}

