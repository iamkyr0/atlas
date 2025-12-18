package resource

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	require.NotNil(t, m)
	require.Greater(t, m.CPUCount, 0)
	require.Greater(t, m.MemoryGB, uint64(0))
	require.Greater(t, m.StorageGB, uint64(0))
}

func TestDetectResources(t *testing.T) {
	m := NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DetectResources(ctx)
	require.NoError(t, err)
	require.NotNil(t, m.NetworkSpeed)
	require.NotNil(t, m.Geolocation)
}

func TestGetResources(t *testing.T) {
	m := NewManager()
	resources := m.GetResources()
	
	require.Contains(t, resources, "cpu_cores")
	require.Contains(t, resources, "memory_gb")
	require.Contains(t, resources, "storage_gb")
	require.Contains(t, resources, "gpu_count")
}

func TestAllocateResources(t *testing.T) {
	m := NewManager()
	
	err := m.AllocateResources("task-1", 2, 4)
	require.NoError(t, err)
	
	require.Equal(t, 2, m.allocatedCPU)
	require.Equal(t, uint64(4), m.allocatedMemory)
	
	err = m.AllocateResources("task-2", m.CPUCount+1, 1)
	require.Error(t, err)
	
	err = m.AllocateResources("task-3", 1, m.MemoryGB+1)
	require.Error(t, err)
}

func TestReleaseResources(t *testing.T) {
	m := NewManager()
	
	m.AllocateResources("task-1", 2, 4)
	require.Equal(t, 2, m.allocatedCPU)
	require.Equal(t, uint64(4), m.allocatedMemory)
	
	err := m.ReleaseResources("task-1")
	require.NoError(t, err)
	
	require.Equal(t, 0, m.allocatedCPU)
	require.Equal(t, uint64(0), m.allocatedMemory)
	
	err = m.ReleaseResources("nonexistent")
	require.NoError(t, err)
}

