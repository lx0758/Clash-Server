package repository

import (
	"clash-server/internal/model"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{db: model.DB}
}

func (r *ConfigRepository) Get(key string) (string, error) {
	var config model.Config
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func (r *ConfigRepository) Set(key, value string) error {
	return r.db.Save(&model.Config{Key: key, Value: value}).Error
}

func (r *ConfigRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&model.Config{}).Error
}

func (r *ConfigRepository) GetMulti(keys []string) (map[string]string, error) {
	var configs []model.Config
	err := r.db.Where("key IN ?", keys).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	return result, nil
}

func (r *ConfigRepository) SetMulti(kv map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range kv {
			if err := tx.Save(&model.Config{Key: key, Value: value}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
