package repository

import (
	"clash-server/internal/model"

	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{db: model.DB}
}

func (r *RuleRepository) Create(rule *model.Rule) error {
	return r.db.Create(rule).Error
}

func (r *RuleRepository) FindAll() ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Order("priority ASC, created_at ASC").Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) FindBySubscriptionID(subscriptionID uint) ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Where("subscription_id = ?", subscriptionID).Order("priority ASC, created_at ASC").Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) FindByID(id uint) (*model.Rule, error) {
	var rule model.Rule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RuleRepository) FindByIDAndSubscriptionID(id, subscriptionID uint) (*model.Rule, error) {
	var rule model.Rule
	err := r.db.Where("id = ? AND subscription_id = ?", id, subscriptionID).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RuleRepository) Update(rule *model.Rule) error {
	return r.db.Save(rule).Error
}

func (r *RuleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Rule{}, id).Error
}

func (r *RuleRepository) DeleteBySubscriptionID(subscriptionID uint) error {
	return r.db.Where("subscription_id = ?", subscriptionID).Delete(&model.Rule{}).Error
}

func (r *RuleRepository) FindByMode(mode model.RuleMode) ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Where("mode = ? AND enabled = ?", mode, true).Order("priority ASC, created_at ASC").Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) FindByModeAndSubscriptionID(mode model.RuleMode, subscriptionID uint) ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Where("mode = ? AND enabled = ? AND subscription_id = ?", mode, true, subscriptionID).Order("priority ASC, created_at ASC").Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) CountBySubscriptionID(subscriptionID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Rule{}).Where("subscription_id = ?", subscriptionID).Count(&count).Error
	return int(count), err
}
