package alert

import (
	"fmt"
	"time"

	"github.com/Mohammad-Ali-Rauf/sentinel.git/pkg/types"
)

type AlertManager struct {
	alerts    []types.Alert
	alertChan chan types.Alert
}

func NewAlertManager() *AlertManager {
	am := &AlertManager{
		alertChan: make(chan types.Alert, 100), // Buffered channel
	}

	// Start alert processor
	go am.processAlerts()

	return am
}

func (am *AlertManager) processAlerts() {
	for alert := range am.alertChan {
		// Store alert
		am.alerts = append(am.alerts, alert)

		// Display alert
		am.displayAlert(alert)
	}
}

func (am *AlertManager) displayAlert(alert types.Alert) {
	timestamp := alert.Time.Format("15:04:05")

	switch alert.Level {
	case "CRITICAL":
		fmt.Printf("🔴 [%s] CRITICAL: %s\n", timestamp, alert.Message)
	case "WARN":
		fmt.Printf("🟡 [%s] WARNING: %s\n", timestamp, alert.Message)
	default:
		fmt.Printf("🔵 [%s] INFO: %s\n", timestamp, alert.Message)
	}
}

func (am *AlertManager) TriggerAlert(level, message string) {
	alert := types.Alert{
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}

	am.alertChan <- alert
}

func (am *AlertManager) GetAlerts() []types.Alert {
	return am.alerts
}

func (am *AlertManager) GetAlertCount() int {
	return len(am.alerts)
}

func (am *AlertManager) Stop() {
	close(am.alertChan)
}
