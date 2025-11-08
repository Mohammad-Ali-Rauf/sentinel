package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Mohammad-Ali-Rauf/sentinel.git/internal/alert"
	"github.com/Mohammad-Ali-Rauf/sentinel.git/pkg/types"
)

type MonitoringService struct {
	config         types.Config
	isRunning      bool
	stopChan       chan bool
	stats          *MonitoringStats
	alertManager   *alert.AlertManager
	knownOpenPorts map[int]bool
}

type MonitoringStats struct {
	ScansCompleted int
	PortsFoundOpen int
	LastScanTime   time.Time
}

func NewMonitoringService(config types.Config, alertManager *alert.AlertManager) *MonitoringService {
	return &MonitoringService{
		config:         config,
		stopChan:       make(chan bool),
		stats:          &MonitoringStats{},
		alertManager:   alertManager,
		knownOpenPorts: make(map[int]bool),
	}
}

func (ms *MonitoringService) TriggerScan() {
	log.Println("🔍 Manual scan triggered by user...")
	ms.performScan()
}

func (ms *MonitoringService) Start() {
	if ms.isRunning {
		log.Println("⚠️ Monitoring service already running")
		return
	}

	ms.isRunning = true
	log.Println("🔍 Starting background monitoring service...")

	// Start the monitoring loop in a goroutine (background)
	go ms.monitoringLoop()
}

func (ms *MonitoringService) Stop() {
	if !ms.isRunning {
		return
	}

	ms.isRunning = false
	ms.stopChan <- true
	log.Println("🛑 Monitoring service stopped")
}

func (ms *MonitoringService) monitoringLoop() {
	ticker := time.NewTicker(time.Duration(ms.config.Thresholds.ScanInterval) * time.Second)
	defer ticker.Stop()

	// Run immediately on start
	ms.performScan()

	for {
		select {
		case <-ticker.C:
			ms.performScan()
		case <-ms.stopChan:
			return
		}
	}
}

func (ms *MonitoringService) performScan() {
	log.Printf("🔍 Scanning %d ports...", len(ms.config.Allow.Ports))

	results := ScanCommonPorts("localhost", ms.config.Allow.Ports)

	openPorts := 0
	for _, result := range results {
		if result.IsOpen {
			openPorts++

			// ALERT: New port detected!
			if !ms.knownOpenPorts[result.Port] {
				ms.knownOpenPorts[result.Port] = true
				ms.alertManager.TriggerAlert("WARN",
					fmt.Sprintf("New open port detected: %d (%s)", result.Port, result.Service))
			}

			log.Printf("✅ Port %d (%s) is OPEN", result.Port, result.Service)
		} else {
			// ALERT: Port that was open is now closed
			if ms.knownOpenPorts[result.Port] {
				delete(ms.knownOpenPorts, result.Port)
				ms.alertManager.TriggerAlert("INFO",
					fmt.Sprintf("Port %d (%s) is now closed", result.Port, result.Service))
			}
		}
	}

	// ALERT: Too many open ports
	if openPorts > ms.config.Thresholds.AlertThreshold {
		ms.alertManager.TriggerAlert("CRITICAL",
			fmt.Sprintf("High number of open ports: %d (threshold: %d)",
				openPorts, ms.config.Thresholds.AlertThreshold))
	}

	ms.stats.ScansCompleted++
	ms.stats.PortsFoundOpen = openPorts
	ms.stats.LastScanTime = time.Now()

	log.Printf("📊 Scan completed: %d/%d ports open", openPorts, len(ms.config.Allow.Ports))
}

func (ms *MonitoringService) GetStats() MonitoringStats {
	if ms.stats == nil {
		return MonitoringStats{}
	}
	return *ms.stats
}

func (ms *MonitoringService) IsRunning() bool {
	return ms.isRunning
}
