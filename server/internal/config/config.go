package config

const MihomoBinary = "clash/mihomo"

type ServerConfig struct {
	Host     string
	Port     int
	Database string
}

type CoreConfig struct {
	Mode      string
	APIHost   string
	APIPort   int
	APISecret string
	MixedPort int
	AllowLan  bool
	LogLevel  string
	IPv6      bool
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:     "0.0.0.0",
		Port:     7000,
		Database: "data.db",
	}
}

func DefaultCoreConfig() *CoreConfig {
	return &CoreConfig{
		Mode:      "rule",
		APIHost:   "127.0.0.1",
		APIPort:   9090,
		MixedPort: 7890,
		AllowLan:  true,
		LogLevel:  "info",
		IPv6:      false,
	}
}
