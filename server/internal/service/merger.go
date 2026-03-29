package service

import (
	"fmt"
	"strings"
	"sync"

	"clash-server/internal/config"
	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/pkg/script"

	"gopkg.in/yaml.v3"
)

type MergerService struct {
	subRepo           *repository.SubscriptionRepository
	customizationRepo *repository.CustomizationRepository
}

var (
	mergerOnce     sync.Once
	mergerInstance *MergerService
)

func GetMergerService() *MergerService {
	mergerOnce.Do(func() {
		mergerInstance = &MergerService{
			subRepo:           repository.NewSubscriptionRepository(),
			customizationRepo: repository.NewCustomizationRepository(),
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

	customization, err := m.customizationRepo.FindBySubscriptionID(subscriptionID)
	if err != nil {
		customization = nil
	}

	if customization != nil {
		baseConfig = m.applyRemoveOperations(baseConfig, customization)
		baseConfig = m.applyInsertOperations(baseConfig, customization)
		baseConfig = m.applyAppendOperations(baseConfig, customization)
	}

	baseConfig = m.applyCoreConfig(baseConfig)

	if customization != nil {
		baseConfig, err = m.applyGlobalOverride(baseConfig, customization)
		if err != nil {
			return nil, fmt.Errorf("全局配置错误: %w", err)
		}
		baseConfig, err = m.applyScript(baseConfig, customization)
		if err != nil {
			return nil, fmt.Errorf("脚本执行错误: %w", err)
		}
	}

	return baseConfig, nil
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

func (m *MergerService) applyRemoveOperations(config map[string]interface{}, customization *model.SubscriptionCustomization) map[string]interface{} {
	config = m.removeProxies(config, customization.ProxyRemove)
	config = m.removeProxyGroups(config, customization.ProxyGroupRemove)
	config = m.removeRules(config, customization.RuleRemove)
	return config
}

func (m *MergerService) applyInsertOperations(config map[string]interface{}, customization *model.SubscriptionCustomization) map[string]interface{} {
	config = m.insertProxies(config, customization.ProxyInsert)
	config = m.insertProxyGroups(config, customization.ProxyGroupInsert)
	config = m.insertRules(config, customization.RuleInsert)
	return config
}

func (m *MergerService) applyAppendOperations(config map[string]interface{}, customization *model.SubscriptionCustomization) map[string]interface{} {
	config = m.appendProxies(config, customization.ProxyAppend)
	config = m.appendProxyGroups(config, customization.ProxyGroupAppend)
	config = m.appendRules(config, customization.RuleAppend)
	return config
}

func (m *MergerService) applyGlobalOverride(config map[string]interface{}, customization *model.SubscriptionCustomization) (map[string]interface{}, error) {
	if customization.GlobalOverride == "" {
		return config, nil
	}
	var override map[string]interface{}
	if err := yaml.Unmarshal([]byte(customization.GlobalOverride), &override); err != nil {
		return nil, fmt.Errorf("全局配置 YAML 格式错误: %w", err)
	}
	for key, value := range override {
		config[key] = value
	}
	return config, nil
}

func (m *MergerService) applyScript(config map[string]interface{}, customization *model.SubscriptionCustomization) (map[string]interface{}, error) {
	if strings.TrimSpace(customization.Script) == "" {
		return config, nil
	}
	engine := script.NewEngine()
	result, err := engine.Execute(customization.Script, config)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MergerService) parseYAMLArray(yamlStr string) []interface{} {
	if yamlStr == "" {
		return nil
	}
	var items []interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &items); err != nil {
		return nil
	}
	return items
}

func (m *MergerService) removeProxies(config map[string]interface{}, yamlStr string) map[string]interface{} {
	names := m.parseYAMLArray(yamlStr)
	if len(names) == 0 {
		return config
	}
	removeSet := make(map[string]bool)
	for _, name := range names {
		if s, ok := name.(string); ok {
			removeSet[s] = true
		}
	}
	proxies, ok := config["proxies"].([]interface{})
	if !ok {
		return config
	}
	var newProxies []interface{}
	for _, p := range proxies {
		if proxy, ok := p.(map[string]interface{}); ok {
			if name, ok := proxy["name"].(string); ok {
				if !removeSet[name] {
					newProxies = append(newProxies, p)
				}
			} else {
				newProxies = append(newProxies, p)
			}
		} else {
			newProxies = append(newProxies, p)
		}
	}
	config["proxies"] = newProxies
	return config
}

func (m *MergerService) insertProxies(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["proxies"].([]interface{})
	config["proxies"] = append(items, existing...)
	return config
}

func (m *MergerService) appendProxies(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["proxies"].([]interface{})
	config["proxies"] = append(existing, items...)
	return config
}

func (m *MergerService) removeProxyGroups(config map[string]interface{}, yamlStr string) map[string]interface{} {
	names := m.parseYAMLArray(yamlStr)
	if len(names) == 0 {
		return config
	}
	removeSet := make(map[string]bool)
	for _, name := range names {
		if s, ok := name.(string); ok {
			removeSet[s] = true
		}
	}
	groups, ok := config["proxy-groups"].([]interface{})
	if !ok {
		return config
	}
	var newGroups []interface{}
	for _, g := range groups {
		if group, ok := g.(map[string]interface{}); ok {
			if name, ok := group["name"].(string); ok {
				if !removeSet[name] {
					newGroups = append(newGroups, g)
				}
			} else {
				newGroups = append(newGroups, g)
			}
		} else {
			newGroups = append(newGroups, g)
		}
	}
	config["proxy-groups"] = newGroups
	return config
}

func (m *MergerService) insertProxyGroups(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["proxy-groups"].([]interface{})
	config["proxy-groups"] = append(items, existing...)
	return config
}

func (m *MergerService) appendProxyGroups(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["proxy-groups"].([]interface{})
	config["proxy-groups"] = append(existing, items...)
	return config
}

func (m *MergerService) removeRules(config map[string]interface{}, yamlStr string) map[string]interface{} {
	rulesToRemove := m.parseYAMLArray(yamlStr)
	if len(rulesToRemove) == 0 {
		return config
	}
	removeSet := make(map[string]bool)
	for _, r := range rulesToRemove {
		if s, ok := r.(string); ok {
			removeSet[s] = true
		}
	}
	rules, ok := config["rules"].([]interface{})
	if !ok {
		return config
	}
	var newRules []interface{}
	for _, r := range rules {
		if s, ok := r.(string); ok {
			if !removeSet[s] {
				newRules = append(newRules, r)
			}
		} else {
			newRules = append(newRules, r)
		}
	}
	config["rules"] = newRules
	return config
}

func (m *MergerService) insertRules(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["rules"].([]interface{})
	config["rules"] = append(items, existing...)
	return config
}

func (m *MergerService) appendRules(config map[string]interface{}, yamlStr string) map[string]interface{} {
	items := m.parseYAMLArray(yamlStr)
	if len(items) == 0 {
		return config
	}
	existing, _ := config["rules"].([]interface{})
	config["rules"] = append(existing, items...)
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
		return fmt.Errorf("配置缺少 proxies 字段")
	}
	proxies, ok := config["proxies"].([]interface{})
	if !ok {
		return fmt.Errorf("proxies 必须是数组")
	}

	availableProxies := make(map[string]bool)
	for _, p := range proxies {
		proxy, ok := p.(map[string]interface{})
		name := extractName(proxy, ok)
		if name == "" {
			return fmt.Errorf("节点缺少 name 字段")
		}
		if !ok {
			return fmt.Errorf("节点 \"%s\" 格式错误", name)
		}
		availableProxies[name] = true
	}

	availableProxyGroups := make(map[string]bool)
	var groups []interface{}
	if g, ok := config["proxy-groups"].([]interface{}); ok {
		groups = g
		for _, g := range groups {
			group, ok := g.(map[string]interface{})
			name := extractName(group, ok)
			if name == "" {
				return fmt.Errorf("代理组缺少 name 字段")
			}
			if !ok {
				return fmt.Errorf("代理组 \"%s\" 格式错误", name)
			}
			availableProxyGroups[name] = true
		}
	}

	for _, g := range groups {
		group, ok := g.(map[string]interface{})
		if !ok {
			continue
		}
		groupName, _ := group["name"].(string)
		if groupProxies, ok := group["proxies"].([]interface{}); ok {
			for _, gp := range groupProxies {
				proxyName, _ := gp.(string)
				if proxyName == "" {
					continue
				}
				if !availableProxies[proxyName] && !availableProxyGroups[proxyName] {
					return fmt.Errorf("代理组 \"%s\" 引用了不存在的节点或策略组: \"%s\"", groupName, proxyName)
				}
			}
		}
	}

	builtinTargets := map[string]bool{
		"DIRECT":      true,
		"REJECT":      true,
		"REJECT-DROP": true,
		"PASS":        true,
		"COMPATIBLE":  true,
	}

	if rules, ok := config["rules"].([]interface{}); ok {
		for _, r := range rules {
			ruleStr, ok := r.(string)
			if !ok {
				return fmt.Errorf("规则格式错误，必须是字符串")
			}
			parts := strings.Split(ruleStr, ",")
			if len(parts) < 2 {
				return fmt.Errorf("规则格式错误: \"%s\"", ruleStr)
			}
			if len(parts) >= 3 {
				target := parts[2]
				if !availableProxies[target] && !availableProxyGroups[target] && !builtinTargets[target] {
					return fmt.Errorf("规则引用了不存在的节点或策略组: \"%s\"", ruleStr)
				}
			}
		}
	}

	return nil
}

func extractName(v map[string]interface{}, ok bool) string {
	if ok {
		if name, ok := v["name"].(string); ok {
			return name
		}
	}
	return ""
}
