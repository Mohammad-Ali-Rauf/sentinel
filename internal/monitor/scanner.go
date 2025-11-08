package monitor

import (
	"fmt"
	"net"
	"time"
)

// PortScanResult represents the result of a port scan
type PortScanResult struct {
	Port    int
	IsOpen  bool
	Service string
}

// ScanPort checks if a specific port is open
func ScanPort(host string, port int) PortScanResult {
	result := PortScanResult{Port: port}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		result.IsOpen = false
		return result
	}

	defer conn.Close()
	result.IsOpen = true
	result.Service = guessService(port)
	return result
}

// ScanCommonPorts scans a list of common ports
func ScanCommonPorts(host string, ports []int) []PortScanResult {
	var results []PortScanResult

	for _, port := range ports {
		result := ScanPort(host, port)
		results = append(results, result)
		time.Sleep(100 * time.Millisecond) // Be nice to the system
	}

	return results
}

func guessService(port int) string {
	services := map[int]string{
		22:   "SSH",
		80:   "HTTP",
		443:  "HTTPS",
		3000: "Node.js",
		5432: "PostgreSQL",
		8080: "Proxy",
		9000: "PHP",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "unknown"
}
