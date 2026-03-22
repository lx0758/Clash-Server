package repository

import (
	"clash-server/internal/model"

	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{db: model.DB}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) FindAll() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Find(&subs).Error
	return subs, err
}

func (r *SubscriptionRepository) FindByID(id uint) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.First(&sub, id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Update(sub *model.Subscription) error {
	return r.db.Save(sub).Error
}

func (r *SubscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Subscription{}, id).Error
}

func (r *SubscriptionRepository) GetActive() (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.Where("active = ?", true).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) SetActive(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Subscription{}).Where("active = ?", true).Update("active", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Subscription{}).Where("id = ?", id).Update("active", true).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SubscriptionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Subscription{}).Count(&count).Error
	return count, err
}
