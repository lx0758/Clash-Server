package config

import (
	"os"
	"strings"

	"clash-server/internal/repository"
)

var (
	serverConfig *ServerConfig
	coreRepo     *repository.ConfigRepository
)

func InitServerConfig() *ServerConfig {
	cfg := DefaultServerConfig()
	if v := os.Getenv("CS_SERVER_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("CS_SERVER_PORT"); v != "" {
		cfg.Port = parseInt(v)
	}
	if v := os.Getenv("CS_SERVER_DATABASE"); v != "" {
		cfg.Database = v
	}
	serverConfig = cfg
	return cfg
}

func InitCoreRepository(repo *repository.ConfigRepository) {
	coreRepo = repo
}

func GetServerConfig() *ServerConfig {
	if serverConfig == nil {
		return DefaultServerConfig()
	}
	return serverConfig
}

func GetCoreConfig() *CoreConfig {
	cfg := DefaultCoreConfig()
	if coreRepo == nil {
		applyCoreEnvVars(cfg)
		return cfg
	}
	if v, err := coreRepo.Get("core:api_host"); err == nil && v != "" {
		cfg.APIHost = v
	} else if env := os.Getenv("CS_CORE_API_HOST"); env != "" {
		cfg.APIHost = env
	}
	if v, err := coreRepo.Get("core:api_port"); err == nil && v != "" {
		cfg.APIPort = parseInt(v)
	} else if env := os.Getenv("CS_CORE_API_PORT"); env != "" {
		cfg.APIPort = parseInt(env)
	}
	if v, err := coreRepo.Get("core:api_secret"); err == nil {
		cfg.APISecret = v
	} else if env := os.Getenv("CS_CORE_API_SECRET"); env != "" {
		cfg.APISecret = env
	}
	if v, err := coreRepo.Get("core:mixed_port"); err == nil && v != "" {
		cfg.MixedPort = parseInt(v)
	} else if env := os.Getenv("CS_CORE_MIXED_PORT"); env != "" {
		cfg.MixedPort = parseInt(env)
	}
	if v, err := coreRepo.Get("core:allow_lan"); err == nil && v != "" {
		cfg.AllowLan = parseBool(v)
	} else if env := os.Getenv("CS_CORE_ALLOW_LAN"); env != "" {
		cfg.AllowLan = parseBool(env)
	}
	if v, err := coreRepo.Get("core:mode"); err == nil && v != "" {
		cfg.Mode = v
	} else if env := os.Getenv("CS_CORE_MODE"); env != "" {
		cfg.Mode = env
	}
	if v, err := coreRepo.Get("core:log_level"); err == nil && v != "" {
		cfg.LogLevel = v
	} else if env := os.Getenv("CS_CORE_LOG_LEVEL"); env != "" {
		cfg.LogLevel = env
	}
	if v, err := coreRepo.Get("core:ipv6"); err == nil && v != "" {
		cfg.IPv6 = parseBool(v)
	} else if env := os.Getenv("CS_CORE_IPV6"); env != "" {
		cfg.IPv6 = parseBool(env)
	}
	return cfg
}

func applyCoreEnvVars(cfg *CoreConfig) {
	if v := os.Getenv("CS_CORE_API_HOST"); v != "" {
		cfg.APIHost = v
	}
	if v := os.Getenv("CS_CORE_API_PORT"); v != "" {
		cfg.APIPort = parseInt(v)
	}
	if v := os.Getenv("CS_CORE_API_SECRET"); v != "" {
		cfg.APISecret = v
	}
	if v := os.Getenv("CS_CORE_MIXED_PORT"); v != "" {
		cfg.MixedPort = parseInt(v)
	}
	if v := os.Getenv("CS_CORE_ALLOW_LAN"); v != "" {
		cfg.AllowLan = parseBool(v)
	}
	if v := os.Getenv("CS_CORE_MODE"); v != "" {
		cfg.Mode = v
	}
	if v := os.Getenv("CS_CORE_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("CS_CORE_IPV6"); v != "" {
		cfg.IPv6 = parseBool(v)
	}
}

func parseInt(s string) int {
	var result int
	for _, c := range strings.TrimSpace(s) {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes" || s == "on"
}
