package conf

type Config struct {
	Generation           int
	WatchNamespace       string
	MyNamespace          string
	MetricsAddr          string
	EnableLeaderElection bool
	ControllerName       string
	LeaderElectionID     string
	SyncDelay            int
	DnsmasqConfDir       string
	DnsmasqOptions       []string
	LogLevel             string
	Development          bool
	CleanupDir           bool
	EnableDHCP           bool
	EnableDNS            bool
}

var instance *Config

func GetConfig() *Config {
	if instance == nil {
		instance = &Config{}
	}
	return instance
}
