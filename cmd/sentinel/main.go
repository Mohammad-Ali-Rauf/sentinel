package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Mohammad-Ali-Rauf/sentinel.git/internal/alert"
	"github.com/Mohammad-Ali-Rauf/sentinel.git/internal/config"
	"github.com/Mohammad-Ali-Rauf/sentinel.git/internal/monitor"
	"github.com/Mohammad-Ali-Rauf/sentinel.git/pkg/types"
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
		IsRunning: false, // Start as NOT running
		StartTime: time.Now(),
	}

	// Create alert manager
	alertManager := alert.NewAlertManager()

	// Create monitoring service BUT DON'T START IT
	monitorService := monitor.NewMonitoringService(cfg, alertManager)

	setupSignalHandling(monitorService, alertManager)

	// Display startup info
	displayStartupInfo(sentinel)

	fmt.Println("✅ Sentinel is ready! Type 'start' to begin monitoring.")

	// Start interactive loop - monitoring is STOPPED by default
	runInteractiveLoop(sentinel, monitorService, alertManager)
}

func displayStartupInfo(s *types.Sentinel) {
	fmt.Printf("🎯 Mode: %s\n", s.Config.Mode)
	fmt.Printf("📊 Allowed Ports: %v\n", s.Config.Allow.Ports)
	fmt.Printf("🚫 Denied Ports: %v\n", s.Config.Deny.Ports)
	fmt.Printf("🌐 Allowed Domains: %v\n", s.Config.Allow.Domains)
	fmt.Printf("⚡ Auto-block: %v\n", s.Config.Thresholds.AutoBlock)
	fmt.Printf("🔍 Scan Interval: %d seconds\n", s.Config.Thresholds.ScanInterval)
	fmt.Printf("🛑 Initial State: STOPPED (type 'start' to begin)\n")
	fmt.Println("----------------------------------------")
}

func runInteractiveLoop(s *types.Sentinel, ms *monitor.MonitoringService, am *alert.AlertManager) {
	for {
		fmt.Println("\nCommands: status, stats, alerts, scan, stop, start, quit")
		fmt.Print("> ")

		var input string
		fmt.Scanln(&input)

		// Convert to lowercase and trim for better UX
		command := strings.ToLower(strings.TrimSpace(input))

		switch command {
		case "status", "s":
			displayStatus(s, ms)
		case "stats", "st":
			displayStats(s, ms)
		case "alerts", "a":
			displayAlerts(am)
		case "scan", "sc":
			fmt.Println("🔍 Triggering immediate scan...")
			ms.TriggerScan()
		case "stop", "x":
			if ms.IsRunning() {
				s.IsRunning = false
				ms.Stop()
				fmt.Println("🛑 Monitoring stopped")
			} else {
				fmt.Println("ℹ️  Monitoring is already stopped")
			}
		case "start", "go":
			if !ms.IsRunning() {
				s.IsRunning = true
				ms.Start()
				fmt.Println("✅ Monitoring started! Scanning for open ports...")
			} else {
				fmt.Println("ℹ️  Monitoring is already running")
			}
		case "quit", "exit", "q":
			ms.Stop()
			am.Stop()
			fmt.Println("👋 Shutting down Sentinel...")
			return
		default:
			fmt.Println("❓ Unknown command. Try: status, stats, alerts, scan, stop, start, quit")
		}
	}
}

func setupSignalHandling(ms *monitor.MonitoringService, am *alert.AlertManager) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n🛑 Received interrupt signal, shutting down gracefully...")
		ms.Stop()
		am.Stop()
		os.Exit(0)
	}()
}

func displayAlerts(am *alert.AlertManager) {
	alerts := am.GetAlerts()
	fmt.Printf("\n📋 Alert History (%d total):\n", len(alerts))

	if len(alerts) == 0 {
		fmt.Println("   No alerts triggered yet")
		return
	}

	// Show last 5 alerts
	start := len(alerts) - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < len(alerts); i++ {
		alert := alerts[i]
		timestamp := alert.Time.Format("15:04:05")

		levelIcon := "🔵"
		switch alert.Level {
		case "CRITICAL":
			levelIcon = "🔴"
		case "WARN":
			levelIcon = "🟡"
		}

		fmt.Printf("   %s [%s] %s: %s\n", levelIcon, timestamp, alert.Level, alert.Message)
	}
}

func displayStatus(s *types.Sentinel, ms *monitor.MonitoringService) {
	status := "🟢 RUNNING"
	if !s.IsRunning {
		status = "🔴 STOPPED"
	}

	monitorStatus := "🔴 STOPPED"
	if ms.IsRunning() {
		monitorStatus = "🟢 RUNNING"
	}

	fmt.Printf("Sentinel Status: %s\n", status)
	fmt.Printf("Monitor Status: %s\n", monitorStatus)
	fmt.Printf("Uptime: %v\n", time.Since(s.StartTime).Round(time.Second))
	fmt.Printf("Mode: %s\n", s.Config.Mode)
}

func displayStats(s *types.Sentinel, ms *monitor.MonitoringService) {
	monitorStats := ms.GetStats()

	stats := types.Stats{
		Uptime:             time.Since(s.StartTime),
		TotalConnections:   monitorStats.ScansCompleted,
		BlockedConnections: 0,
		AlertsTriggered:    0,
		ActiveMonitors:     1,
	}

	fmt.Println("📊 Current Statistics:")
	fmt.Printf("   Uptime: %v\n", stats.Uptime.Round(time.Second))
	fmt.Printf("   Scans Completed: %d\n", monitorStats.ScansCompleted)
	fmt.Printf("   Open Ports Found: %d\n", monitorStats.PortsFoundOpen)
	fmt.Printf("   Last Scan: %v\n", monitorStats.LastScanTime.Format("15:04:05"))
	fmt.Printf("   Alerts Triggered: %d\n", stats.AlertsTriggered)
}
