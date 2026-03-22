package repository

import (
	"clash-server/internal/model"

	"gorm.io/gorm"
)

type ScriptRepository struct {
	db *gorm.DB
}

func NewScriptRepository() *ScriptRepository {
	return &ScriptRepository{db: model.DB}
}

func (r *ScriptRepository) Create(script *model.Script) error {
	return r.db.Create(script).Error
}

func (r *ScriptRepository) FindAll() ([]model.Script, error) {
	var scripts []model.Script
	err := r.db.Order("created_at DESC").Find(&scripts).Error
	return scripts, err
}

func (r *ScriptRepository) FindBySubscriptionID(subscriptionID uint) ([]model.Script, error) {
	var scripts []model.Script
	err := r.db.Where("subscription_id = ?", subscriptionID).Order("created_at DESC").Find(&scripts).Error
	return scripts, err
}

func (r *ScriptRepository) FindByID(id uint) (*model.Script, error) {
	var script model.Script
	err := r.db.First(&script, id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *ScriptRepository) FindByIDAndSubscriptionID(id, subscriptionID uint) (*model.Script, error) {
	var script model.Script
	err := r.db.Where("id = ? AND subscription_id = ?", id, subscriptionID).First(&script).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *ScriptRepository) FindEnabled() ([]model.Script, error) {
	var scripts []model.Script
	err := r.db.Where("enabled = ?", true).Order("created_at ASC").Find(&scripts).Error
	return scripts, err
}

func (r *ScriptRepository) FindEnabledBySubscriptionID(subscriptionID uint) ([]model.Script, error) {
	var scripts []model.Script
	err := r.db.Where("enabled = ? AND subscription_id = ?", true, subscriptionID).Order("created_at ASC").Find(&scripts).Error
	return scripts, err
}

func (r *ScriptRepository) Update(script *model.Script) error {
	return r.db.Save(script).Error
}

func (r *ScriptRepository) Delete(id uint) error {
	return r.db.Delete(&model.Script{}, id).Error
}

func (r *ScriptRepository) DeleteBySubscriptionID(subscriptionID uint) error {
	return r.db.Where("subscription_id = ?", subscriptionID).Delete(&model.Script{}).Error
}

func (r *ScriptRepository) CountBySubscriptionID(subscriptionID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Script{}).Where("subscription_id = ?", subscriptionID).Count(&count).Error
	return int(count), err
}
