package resource

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/atlas/node/network"
)

type Manager struct {
	CPUCount    int
	GPUs        []GPU
	MemoryGB    uint64
	StorageGB   uint64
	NetworkSpeed *network.SpeedTestResult
	Geolocation *network.Geolocation
	allocations map[string]*ResourceAllocation
	allocatedCPU int
	allocatedMemory uint64
}

type GPU struct {
	ID          string
	MemoryGB    uint64
	Utilization float64
}

func NewManager() *Manager {
	m := &Manager{
		CPUCount: runtime.NumCPU(),
		GPUs:     []GPU{},
	}
	
	// Auto-detect all resources
	m.detectMemory()
	m.detectStorage()
	m.detectGPUs()
	
	// Network and geolocation are async and can be done later
	// They're expensive operations, so we'll do them on demand
	
	return m
}

// DetectResources performs full resource detection including network tests
func (m *Manager) DetectResources(ctx context.Context) error {
	// Detect network speed
	speedResult, err := network.SpeedTest(ctx)
	if err != nil {
		return fmt.Errorf("speed test failed: %w", err)
	}
	m.NetworkSpeed = speedResult
	
	// Detect geolocation
	geo, err := network.GetGeolocation(ctx)
	if err != nil {
		return fmt.Errorf("geolocation failed: %w", err)
	}
	m.Geolocation = geo
	
	return nil
}

func (m *Manager) detectMemory() {
	// Get system memory
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err == nil {
		// Convert to GB (info.Totalram is in bytes)
		m.MemoryGB = info.Totalram / (1024 * 1024 * 1024)
	} else {
		// Fallback: use runtime memory stats (less accurate)
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		// This gives allocated memory, not total system memory
		// For a better estimate, we'd need platform-specific code
		m.MemoryGB = 8 // Default fallback
	}
}

func (m *Manager) detectStorage() {
	// Get disk space using syscall
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err == nil {
		// Calculate total space in GB
		totalBytes := stat.Blocks * uint64(stat.Bsize)
		m.StorageGB = totalBytes / (1024 * 1024 * 1024)
	} else {
		// Fallback: try using df command
		cmd := exec.Command("df", "-BG", "/")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) > 1 {
				fields := strings.Fields(lines[1])
				if len(fields) > 1 {
					// Remove 'G' suffix and parse
					sizeStr := strings.TrimSuffix(fields[1], "G")
					if size, err := strconv.ParseUint(sizeStr, 10, 64); err == nil {
						m.StorageGB = size
					}
				}
			}
		}
	}
}

func (m *Manager) detectGPUs() {
	// Try NVIDIA GPUs first (nvidia-smi)
	m.detectNvidiaGPUs()
	
	// Could add other GPU detection here (AMD, Intel, etc.)
}

func (m *Manager) detectNvidiaGPUs() {
	// Check if nvidia-smi is available
	cmd := exec.Command("nvidia-smi", "--query-gpu=index,name,memory.total", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		// No NVIDIA GPUs or nvidia-smi not available
		return
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		fields := strings.Split(line, ",")
		if len(fields) >= 3 {
			gpuID := strings.TrimSpace(fields[0])
			memoryMBStr := strings.TrimSpace(fields[2])
			
			memoryMB, err := strconv.ParseUint(memoryMBStr, 10, 64)
			if err == nil {
				memoryGB := memoryMB / 1024
				m.GPUs = append(m.GPUs, GPU{
					ID:       gpuID,
					MemoryGB: memoryGB,
					Utilization: 0.0, // Would need to query separately
				})
			}
		}
	}
}

// GetResources returns resource information as a map
func (m *Manager) GetResources() map[string]interface{} {
	resources := map[string]interface{}{
		"cpu_cores": m.CPUCount,
		"memory_gb": m.MemoryGB,
		"storage_gb": m.StorageGB,
		"gpu_count":  len(m.GPUs),
	}
	
	if m.NetworkSpeed != nil {
		resources["network_download_mbps"] = m.NetworkSpeed.DownloadSpeedMbps
		resources["network_upload_mbps"] = m.NetworkSpeed.UploadSpeedMbps
		resources["network_latency_ms"] = m.NetworkSpeed.LatencyMs
	}
	
	if m.Geolocation != nil {
		resources["region"] = m.Geolocation.Region
		resources["country"] = m.Geolocation.Country
		resources["ip"] = m.Geolocation.IP
	}
	
	return resources
}

type ResourceAllocation struct {
	TaskID string
	CPU    int
	Memory uint64
}

func NewManager() *Manager {
	m := &Manager{
		CPUCount: runtime.NumCPU(),
		GPUs:     []GPU{},
		allocations: make(map[string]*ResourceAllocation),
	}
	
	m.detectMemory()
	m.detectStorage()
	m.detectGPUs()
	
	return m
}

func (m *Manager) AllocateResources(taskID string, cpu int, memory uint64) error {
	if cpu > m.CPUCount {
		return fmt.Errorf("insufficient CPU: requested %d, available %d", cpu, m.CPUCount)
	}
	
	if memory > m.MemoryGB {
		return fmt.Errorf("insufficient memory: requested %d GB, available %d GB", memory, m.MemoryGB)
	}

	availableCPU := m.CPUCount - m.allocatedCPU
	if cpu > availableCPU {
		return fmt.Errorf("insufficient CPU: requested %d, available %d (total %d, allocated %d)", cpu, availableCPU, m.CPUCount, m.allocatedCPU)
	}

	availableMemory := m.MemoryGB - m.allocatedMemory
	if memory > availableMemory {
		return fmt.Errorf("insufficient memory: requested %d GB, available %d GB (total %d GB, allocated %d GB)", memory, availableMemory, m.MemoryGB, m.allocatedMemory)
	}

	m.allocations[taskID] = &ResourceAllocation{
		TaskID: taskID,
		CPU:    cpu,
		Memory: memory,
	}
	m.allocatedCPU += cpu
	m.allocatedMemory += memory

	return nil
}

func (m *Manager) ReleaseResources(taskID string) error {
	allocation, exists := m.allocations[taskID]
	if !exists {
		return nil
	}

	m.allocatedCPU -= allocation.CPU
	m.allocatedMemory -= allocation.Memory
	delete(m.allocations, taskID)

	return nil
}

