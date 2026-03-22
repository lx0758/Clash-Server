package service

import (
	"errors"

	"clash-server/internal/model"
	"clash-server/internal/repository"
)

var ValidRuleTypes = map[string]bool{
	"DOMAIN":         true,
	"DOMAIN-SUFFIX":  true,
	"DOMAIN-KEYWORD": true,
	"IP-CIDR":        true,
	"IP-CIDR6":       true,
	"SRC-IP-CIDR":    true,
	"GEOIP":          true,
	"GEOSITE":        true,
	"SRC-GEOIP":      true,
	"DST-PORT":       true,
	"SRC-PORT":       true,
	"PROCESS-NAME":   true,
	"RULE-SET":       true,
	"MATCH":          true,
}

type RuleService struct {
	repo *repository.RuleRepository
}

func NewRuleService() *RuleService {
	return &RuleService{repo: repository.NewRuleRepository()}
}

func (s *RuleService) List(subscriptionID uint) ([]model.Rule, error) {
	return s.repo.FindBySubscriptionID(subscriptionID)
}

func (s *RuleService) Get(id uint) (*model.Rule, error) {
	return s.repo.FindByID(id)
}

func (s *RuleService) GetBySubscription(id, subscriptionID uint) (*model.Rule, error) {
	return s.repo.FindByIDAndSubscriptionID(id, subscriptionID)
}

func (s *RuleService) Create(rule *model.Rule) error {
	if !ValidRuleTypes[rule.Type] {
		return errors.New("invalid rule type")
	}
	return s.repo.Create(rule)
}

func (s *RuleService) Update(rule *model.Rule) error {
	if !ValidRuleTypes[rule.Type] {
		return errors.New("invalid rule type")
	}
	return s.repo.Update(rule)
}

func (s *RuleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *RuleService) DeleteBySubscription(subscriptionID uint) error {
	return s.repo.DeleteBySubscriptionID(subscriptionID)
}

func (s *RuleService) GetInsertRules(subscriptionID uint) ([]model.Rule, error) {
	return s.repo.FindByModeAndSubscriptionID(model.RuleModeInsert, subscriptionID)
}

func (s *RuleService) GetAppendRules(subscriptionID uint) ([]model.Rule, error) {
	return s.repo.FindByModeAndSubscriptionID(model.RuleModeAppend, subscriptionID)
}

func (s *RuleService) ValidateRuleType(ruleType string) bool {
	return ValidRuleTypes[ruleType]
}
