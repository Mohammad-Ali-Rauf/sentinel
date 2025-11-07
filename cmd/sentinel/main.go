package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Mohammad-Ali-Rauf/sentinel.git/pkg/types"

	"github.com/Mohammad-Ali-Rauf/sentinel.git/internal/config"
)

func main() {
	fmt.Println("🛡️  Sentinel - Network Intrusion Detector")
	fmt.Println("==========================================")

	// Load configuration
	cfg, err := config.LoadConfig("configs/sentinel.toml")
	if err != nil {
		log.Printf("⚠️  Could not load config: %v", err)
		log.Println("🔄 Using default configuration...")
		cfg = types.Config{Mode: types.ModeDev}
		cfg.ApplyPreset()
	}

	// Create Sentinel instance
	sentinel := &types.Sentinel{
		Config:    cfg,
		IsRunning: false,
		StartTime: time.Now(),
	}

	// Display startup info
	displayStartupInfo(sentinel)

	// Start basic monitoring
	sentinel.IsRunning = true
	fmt.Println("✅ Sentinel is now running!")

	// Simple interactive loop
	runInteractiveLoop(sentinel)
}

func displayStartupInfo(s *types.Sentinel) {
	fmt.Printf("🎯 Mode: %s\n", s.Config.Mode)
	fmt.Printf("📊 Allowed Ports: %v\n", s.Config.Allow.Ports)
	fmt.Printf("🚫 Denied Ports: %v\n", s.Config.Deny.Ports)
	fmt.Printf("🌐 Allowed Domains: %v\n", s.Config.Allow.Domains)
	fmt.Printf("⚡ Auto-block: %v\n", s.Config.Thresholds.AutoBlock)
	fmt.Printf("🔍 Scan Interval: %d seconds\n", s.Config.Thresholds.ScanInterval)
	fmt.Println("----------------------------------------")
}

func runInteractiveLoop(s *types.Sentinel) {
	for {
		fmt.Println("\nCommands: status, stats, stop, quit")
		fmt.Print("> ")

		var input string
		fmt.Scanln(&input)

		switch input {
		case "status":
			displayStatus(s)
		case "stats":
			displayStats(s)
		case "stop":
			s.IsRunning = false
			fmt.Println("🛑 Monitoring stopped")
		case "quit", "exit":
			fmt.Println("👋 Shutting down Sentinel...")
			return
		default:
			fmt.Println("❓ Unknown command. Try: status, stats, stop, quit")
		}
	}
}

func displayStatus(s *types.Sentinel) {
	status := "🟢 RUNNING"
	if !s.IsRunning {
		status = "🔴 STOPPED"
	}
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Uptime: %v\n", time.Since(s.StartTime).Round(time.Second))
	fmt.Printf("Mode: %s\n", s.Config.Mode)
}

func displayStats(s *types.Sentinel) {
	stats := types.Stats{
		Uptime:             time.Since(s.StartTime),
		TotalConnections:   0, // We'll implement these later
		BlockedConnections: 0,
		AlertsTriggered:    0,
		ActiveMonitors:     1,
	}

	fmt.Println("📊 Current Statistics:")
	fmt.Printf("   Uptime: %v\n", stats.Uptime.Round(time.Second))
	fmt.Printf("   Total Connections: %d\n", stats.TotalConnections)
	fmt.Printf("   Blocked Connections: %d\n", stats.BlockedConnections)
	fmt.Printf("   Alerts Triggered: %d\n", stats.AlertsTriggered)
	fmt.Printf("   Active Monitors: %d\n", stats.ActiveMonitors)
}
