package types

import "time"

// Preset modes
const (
	ModeDev      = "dev"
	ModeStrict   = "strict"
	ModePassive  = "passive"
	ModeHoneypot = "honeypot"
)

// Config represents the main configuration structure
type Config struct {
	Mode       string     `toml:"mode"`
	Allow      Allow      `toml:"allow"`
	Deny       Deny       `toml:"deny"`
	Thresholds Thresholds `toml:"thresholds"`
}

// ApplyPreset auto-fills configuration based on the selected preset
func (c *Config) ApplyPreset() {
	switch c.Mode {
	case ModeDev:
		c.applyDevPreset()
	case ModeStrict:
		c.applyStrictPreset()
	case ModePassive:
		c.applyPassivePreset()
	case ModeHoneypot:
		c.applyHoneypotPreset()
	}
	// If mode is custom or unknown, don't auto-fill
}

func (c *Config) applyDevPreset() {
	c.Allow.Ports = []int{22, 80, 443, 3000, 5432, 8080} // Common dev ports
	c.Allow.Domains = []string{"github.com", "docker.io", "npmjs.com"}
	c.Deny.Ports = []int{23, 135, 445, 3389} // Commonly exploited ports
	c.Thresholds = Thresholds{
		MaxConnectionsPerMinute: 500,
		AlertOnNewListener:      true,
		ScanInterval:            30,
		AlertThreshold:          10,
		AutoBlock:               false, // Don't block in dev mode
	}
}

func (c *Config) applyStrictPreset() {
	c.Allow.Ports = []int{22, 80, 443} // Only essential ports
	c.Allow.Domains = []string{}       // No domain whitelist
	c.Deny.Ports = []int{23, 135, 445, 3389, 5900, 6379}
	c.Thresholds = Thresholds{
		MaxConnectionsPerMinute: 50,
		AlertOnNewListener:      true,
		ScanInterval:            10,
		AlertThreshold:          1,    // Alert on first suspicious activity
		AutoBlock:               true, // Auto-block in strict mode
	}
}

func (c *Config) applyPassivePreset() {
	c.Allow.Ports = []int{22, 80, 443, 3000, 5432, 8080, 9000}
	c.Allow.Domains = []string{"github.com", "docker.io"}
	c.Deny.Ports = []int{}
	c.Thresholds = Thresholds{
		MaxConnectionsPerMinute: 1000,
		AlertOnNewListener:      true,
		ScanInterval:            60,
		AlertThreshold:          50,    // Only alert on high activity
		AutoBlock:               false, // Never auto-block
	}
}

func (c *Config) applyHoneypotPreset() {
	c.Allow.Ports = []int{} // Allow nothing
	c.Allow.Domains = []string{}
	c.Deny.Ports = []int{22, 80, 443, 3389, 5900} // Common attack ports
	c.Thresholds = Thresholds{
		MaxConnectionsPerMinute: 5, // Very low threshold
		AlertOnNewListener:      true,
		ScanInterval:            5,     // Frequent scanning
		AlertThreshold:          1,     // Alert on everything
		AutoBlock:               false, // Don't block, just monitor
	}
}

type Allow struct {
	Domains []string `toml:"domains"`
	Ports   []int    `toml:"ports"`
}

type Deny struct {
	Ports []int `toml:"ports"`
}

type Thresholds struct {
	MaxConnectionsPerMinute int  `toml:"max_connections_per_minute"`
	AlertOnNewListener      bool `toml:"alert_on_new_listener"`
	ScanInterval            int  `toml:"scan_interval_seconds"`
	AlertThreshold          int  `toml:"alert_threshold"`
	AutoBlock               bool `toml:"auto_block"`
}

// Alert represents a security alert
type Alert struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// Sentinel represents the main application state
type Sentinel struct {
	Config    Config
	IsRunning bool
	StartTime time.Time
	Uptime    time.Duration
}

// Stats represents runtime statistics
type Stats struct {
	Uptime             time.Duration
	TotalConnections   int
	BlockedConnections int
	AlertsTriggered    int
	ActiveMonitors     int
}
