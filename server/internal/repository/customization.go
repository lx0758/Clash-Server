package repository

import (
	"clash-server/internal/model"

	"gorm.io/gorm"
)

type CustomizationRepository struct {
	db *gorm.DB
}

func NewCustomizationRepository() *CustomizationRepository {
	return &CustomizationRepository{db: model.DB}
}

func (r *CustomizationRepository) FindBySubscriptionID(subscriptionID uint) (*model.SubscriptionCustomization, error) {
	var c model.SubscriptionCustomization
	err := r.db.Where("subscription_id = ?", subscriptionID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomizationRepository) Create(c *model.SubscriptionCustomization) error {
	return r.db.Create(c).Error
}

func (r *CustomizationRepository) Update(c *model.SubscriptionCustomization) error {
	return r.db.Save(c).Error
}

func (r *CustomizationRepository) DeleteBySubscriptionID(subscriptionID uint) error {
	return r.db.Where("subscription_id = ?", subscriptionID).Delete(&model.SubscriptionCustomization{}).Error
}

func (r *CustomizationRepository) Upsert(c *model.SubscriptionCustomization) error {
	var existing model.SubscriptionCustomization
	err := r.db.Where("subscription_id = ?", c.SubscriptionID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(c).Error
	}
	if err != nil {
		return err
	}
	c.ID = existing.ID
	return r.db.Save(c).Error
}
