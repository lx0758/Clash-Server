package service

import (
	"fmt"
	"sync"

	"clash-server/internal/config"
	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/pkg/script"

	"gopkg.in/yaml.v3"
)

type MergerService struct {
	subRepo     *repository.SubscriptionRepository
	ruleService *RuleService
	scriptRepo  *repository.ScriptRepository
}

var (
	mergerOnce     sync.Once
	mergerInstance *MergerService
)

func GetMergerService() *MergerService {
	mergerOnce.Do(func() {
		mergerInstance = &MergerService{
			subRepo:     repository.NewSubscriptionRepository(),
			ruleService: NewRuleService(),
			scriptRepo:  repository.NewScriptRepository(),
		}
	})
	return mergerInstance
}

func NewMergerService() *MergerService {
	return GetMergerService()
}

func (m *MergerService) Merge() (map[string]interface{}, error) {
	sub, err := m.subRepo.GetActive()
	if err != nil || sub == nil {
		return m.GetMinimalConfig(), nil
	}
	return m.MergeForSubscription(sub.ID)
}

func (m *MergerService) MergeForSubscription(subscriptionID uint) (map[string]interface{}, error) {
	sub, err := m.subRepo.FindByID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	var baseConfig map[string]interface{}
	if sub.Content != "" {
		if err := yaml.Unmarshal([]byte(sub.Content), &baseConfig); err != nil {
			return nil, fmt.Errorf("failed to parse subscription: %w", err)
		}
	} else {
		baseConfig = make(map[string]interface{})
	}
	insertRules, err := m.ruleService.GetInsertRules(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get insert rules: %w", err)
	}
	appendRules, err := m.ruleService.GetAppendRules(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get append rules: %w", err)
	}
	baseConfig = m.applyInsertRules(baseConfig, insertRules)
	baseConfig = m.applyCoreConfig(baseConfig)
	baseConfig = m.applyAppendRules(baseConfig, appendRules)
	scripts, err := m.scriptRepo.FindEnabledBySubscriptionID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scripts: %w", err)
	}
	result := baseConfig
	engine := script.NewEngine()
	for _, scr := range scripts {
		result, err = engine.Execute(scr.Content, result)
		if err != nil {
			return nil, fmt.Errorf("failed to execute script %s: %w", scr.Name, err)
		}
	}
	return result, nil
}

func (m *MergerService) GetMinimalConfig() map[string]interface{} {
	coreCfg := config.GetCoreConfig()
	minimal := map[string]interface{}{
		"mixed-port":          coreCfg.MixedPort,
		"allow-lan":           coreCfg.AllowLan,
		"mode":                coreCfg.Mode,
		"log-level":           coreCfg.LogLevel,
		"ipv6":                coreCfg.IPv6,
		"external-controller": fmt.Sprintf("%s:%d", coreCfg.APIHost, coreCfg.APIPort),
	}
	if coreCfg.APISecret != "" {
		minimal["secret"] = coreCfg.APISecret
	}
	return minimal
}

func (m *MergerService) applyCoreConfig(cfg map[string]interface{}) map[string]interface{} {
	coreCfg := config.GetCoreConfig()
	cfg["mixed-port"] = coreCfg.MixedPort
	cfg["allow-lan"] = coreCfg.AllowLan
	cfg["mode"] = coreCfg.Mode
	cfg["log-level"] = coreCfg.LogLevel
	cfg["ipv6"] = coreCfg.IPv6
	cfg["external-controller"] = fmt.Sprintf("%s:%d", coreCfg.APIHost, coreCfg.APIPort)
	if coreCfg.APISecret != "" {
		cfg["secret"] = coreCfg.APISecret
	}
	return cfg
}

func (m *MergerService) applyInsertRules(config map[string]interface{}, rules []model.Rule) map[string]interface{} {
	if len(rules) == 0 {
		return config
	}
	ruleStrings := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleStrings = append(ruleStrings, fmt.Sprintf("%s,%s,%s", rule.Type, rule.Payload, rule.Proxy))
	}
	if existingRules, ok := config["rules"].([]interface{}); ok {
		config["rules"] = append(ruleStrings, existingRules...)
	} else {
		config["rules"] = ruleStrings
	}
	return config
}

func (m *MergerService) applyAppendRules(config map[string]interface{}, rules []model.Rule) map[string]interface{} {
	if len(rules) == 0 {
		return config
	}
	ruleStrings := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleStrings = append(ruleStrings, fmt.Sprintf("%s,%s,%s", rule.Type, rule.Payload, rule.Proxy))
	}
	if existingRules, ok := config["rules"].([]interface{}); ok {
		config["rules"] = append(existingRules, ruleStrings...)
	} else {
		config["rules"] = ruleStrings
	}
	return config
}

func (m *MergerService) applyPanelConfig(config map[string]interface{}, panelConfig map[string]interface{}) map[string]interface{} {
	for key, value := range panelConfig {
		config[key] = value
	}
	return config
}

func (m *MergerService) GenerateYAML(config map[string]interface{}) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to generate yaml: %w", err)
	}
	return string(data), nil
}

func (m *MergerService) Validate(config map[string]interface{}) error {
	if _, ok := config["proxies"]; !ok {
		return fmt.Errorf("missing proxies in config")
	}
	return nil
}
