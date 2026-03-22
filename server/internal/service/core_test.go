package service

import (
	"reflect"
	"testing"
)

func TestNeedsRestart(t *testing.T) {
	tests := []struct {
		name        string
		oldConfig   map[string]interface{}
		newConfig   map[string]interface{}
		needRestart bool
	}{
		{
			name:        "no change",
			oldConfig:   map[string]interface{}{"mixed-port": 7890, "proxies": []string{"a"}},
			newConfig:   map[string]interface{}{"mixed-port": 7890, "proxies": []string{"b"}},
			needRestart: false,
		},
		{
			name:        "mixed-port changed",
			oldConfig:   map[string]interface{}{"mixed-port": 7890},
			newConfig:   map[string]interface{}{"mixed-port": 7891},
			needRestart: true,
		},
		{
			name:        "external-controller changed",
			oldConfig:   map[string]interface{}{"external-controller": "127.0.0.1:9090"},
			newConfig:   map[string]interface{}{"external-controller": "127.0.0.1:9091"},
			needRestart: true,
		},
		{
			name:        "secret changed",
			oldConfig:   map[string]interface{}{"secret": "old"},
			newConfig:   map[string]interface{}{"secret": "new"},
			needRestart: true,
		},
		{
			name:        "ipv6 changed",
			oldConfig:   map[string]interface{}{"ipv6": false},
			newConfig:   map[string]interface{}{"ipv6": true},
			needRestart: true,
		},
		{
			name:        "proxies changed - no restart needed",
			oldConfig:   map[string]interface{}{"proxies": []string{"proxy1"}},
			newConfig:   map[string]interface{}{"proxies": []string{"proxy1", "proxy2"}},
			needRestart: false,
		},
		{
			name:        "rules changed - no restart needed",
			oldConfig:   map[string]interface{}{"rules": []string{"DOMAIN,test.com,DIRECT"}},
			newConfig:   map[string]interface{}{"rules": []string{"DOMAIN,test.com,PROXY"}},
			needRestart: false,
		},
		{
			name:        "tun changed",
			oldConfig:   map[string]interface{}{"tun": map[string]interface{}{"enable": false}},
			newConfig:   map[string]interface{}{"tun": map[string]interface{}{"enable": true}},
			needRestart: true,
		},
	}

	cs := &CoreService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cs.needsRestart(tt.oldConfig, tt.newConfig)
			if result != tt.needRestart {
				t.Errorf("needsRestart() = %v, want %v", result, tt.needRestart)
			}
		})
	}
}

func TestRestartRequiredKeys(t *testing.T) {
	expectedKeys := []string{
		"mixed-port",
		"external-controller",
		"secret",
		"bind-address",
		"tun",
		"interface-name",
		"routing-mark",
		"ipv6",
	}

	if !reflect.DeepEqual(RestartRequiredKeys, expectedKeys) {
		t.Errorf("RestartRequiredKeys = %v, want %v", RestartRequiredKeys, expectedKeys)
	}
}
