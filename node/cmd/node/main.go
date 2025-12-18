package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atlas/node/executor"
	"github.com/atlas/node/health"
	"github.com/atlas/node/resource"
	"github.com/spf13/cobra"
)

var (
	chainRPCURL string
	ipfsAPIURL  string
	nodeID      string
	nodeAddress string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "atlas-node",
		Short: "Atlas Compute Node",
		Long:  "Atlas Decentralized AI Platform - Compute Node",
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&chainRPCURL, "chain-rpc", "http://localhost:26657", "Chain RPC URL")
	rootCmd.PersistentFlags().StringVar(&ipfsAPIURL, "ipfs-api", "/ip4/127.0.0.1/tcp/5001", "IPFS API URL")
	rootCmd.PersistentFlags().StringVar(&nodeID, "node-id", "", "Node ID")
	rootCmd.PersistentFlags().StringVar(&nodeAddress, "address", "", "Node wallet address")

	// Start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Initialize components with auto-detection
			resourceManager := resource.NewManager()
			
			// Perform full resource detection
			fmt.Println("Detecting system resources...")
			if err := resourceManager.DetectResources(ctx); err != nil {
				fmt.Printf("Warning: Resource detection failed: %v\n", err)
			}
			
			// Print detected resources
			resources := resourceManager.GetResources()
			fmt.Println("Detected resources:")
			resourceJSON, _ := json.MarshalIndent(resources, "", "  ")
			fmt.Println(string(resourceJSON))

			executor := executor.NewExecutor(resourceManager)
			healthMonitor := health.NewMonitor()

			// Start services
			fmt.Println("Starting node services...")
			go healthMonitor.Start(ctx)
			go executor.Start(ctx)

			// Wait for interrupt
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			<-sigChan

			fmt.Println("Shutting down node...")
			return nil
		},
	}

	// Status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show node status",
		RunE: func(cmd *cobra.Command, args []string) error {
			resourceManager := resource.NewManager()
			
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			if err := resourceManager.DetectResources(ctx); err != nil {
				return fmt.Errorf("resource detection failed: %w", err)
			}
			
			resources := resourceManager.GetResources()
			
			fmt.Println("Node Status:")
			fmt.Println("============")
			fmt.Printf("CPU Cores: %d\n", resourceManager.CPUCount)
			fmt.Printf("Memory: %d GB\n", resourceManager.MemoryGB)
			fmt.Printf("Storage: %d GB\n", resourceManager.StorageGB)
			fmt.Printf("GPUs: %d\n", len(resourceManager.GPUs))
			
			if resourceManager.NetworkSpeed != nil {
				fmt.Printf("\nNetwork:\n")
				fmt.Printf("  Download: %.2f Mbps\n", resourceManager.NetworkSpeed.DownloadSpeedMbps)
				fmt.Printf("  Upload: %.2f Mbps\n", resourceManager.NetworkSpeed.UploadSpeedMbps)
				fmt.Printf("  Latency: %.2f ms\n", resourceManager.NetworkSpeed.LatencyMs)
			}
			
			if resourceManager.Geolocation != nil {
				fmt.Printf("\nLocation:\n")
				fmt.Printf("  IP: %s\n", resourceManager.Geolocation.IP)
				fmt.Printf("  Country: %s\n", resourceManager.Geolocation.Country)
				fmt.Printf("  Region: %s\n", resourceManager.Geolocation.Region)
				fmt.Printf("  City: %s\n", resourceManager.Geolocation.City)
			}
			
			return nil
		},
	}

	// Register command
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register node on blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if nodeID == "" {
				return fmt.Errorf("node-id is required")
			}
			if nodeAddress == "" {
				return fmt.Errorf("address is required")
			}

			resourceManager := resource.NewManager()
			
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			if err := resourceManager.DetectResources(ctx); err != nil {
				return fmt.Errorf("resource detection failed: %w", err)
			}

			resources := resourceManager.GetResources()
			
			fmt.Printf("Registering node %s...\n", nodeID)
			fmt.Printf("Resources:\n")
			resourceJSON, _ := json.MarshalIndent(resources, "", "  ")
			fmt.Println(string(resourceJSON))
			
			cpuCores, _ := resources["cpu_cores"].(int)
			if cpuCores == 0 {
				if c, ok := resources["cpu_cores"].(float64); ok {
					cpuCores = int(c)
				}
			}
			gpuCount, _ := resources["gpu_count"].(int)
			if gpuCount == 0 {
				if g, ok := resources["gpu_count"].(float64); ok {
					gpuCount = int(g)
				}
			}
			memoryGB, _ := resources["memory_gb"].(uint64)
			if memoryGB == 0 {
				if m, ok := resources["memory_gb"].(float64); ok {
					memoryGB = uint64(m)
				}
			}
			storageGB, _ := resources["storage_gb"].(uint64)
			if storageGB == 0 {
				if s, ok := resources["storage_gb"].(float64); ok {
					storageGB = uint64(s)
				}
			}

			err := registerNodeOnBlockchain(chainRPCURL, nodeID, nodeAddress, cpuCores, gpuCount, int(memoryGB), int(storageGB))
			if err != nil {
				return fmt.Errorf("blockchain registration failed: %w\nNote: Use 'atlasd tx compute register-node ...' as fallback", err)
			}

			fmt.Println("Node registered successfully on blockchain")
			return nil
		},
	}

	// Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show or update node configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Node Configuration:")
			fmt.Println("===================")
			fmt.Printf("Chain RPC URL: %s\n", chainRPCURL)
			fmt.Printf("IPFS API URL: %s\n", ipfsAPIURL)
			fmt.Printf("Node ID: %s\n", nodeID)
			fmt.Printf("Node Address: %s\n", nodeAddress)
			return nil
		},
	}

	rootCmd.AddCommand(startCmd, statusCmd, registerCmd, configCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func registerNodeOnBlockchain(rpcURL string, nodeID string, address string, cpuCores int, gpuCount int, memoryGB int, storageGB int) error {
	return fmt.Errorf("blockchain client not fully implemented: requires gRPC client with protobuf stubs generated from chain proto files. Use 'atlasd tx compute register-node --node-id %s --address %s --cpu-cores %d --gpu-count %d --memory-gb %d --storage-gb %d' as alternative", nodeID, address, cpuCores, gpuCount, memoryGB, storageGB)
}

