package service

import (
	"sync"

	"clash-server/internal/model"
	"clash-server/internal/repository"
)

type CustomizationService struct {
	repo *repository.CustomizationRepository
}

var (
	customizationOnce     sync.Once
	customizationInstance *CustomizationService
)

func GetCustomizationService() *CustomizationService {
	customizationOnce.Do(func() {
		customizationInstance = &CustomizationService{
			repo: repository.NewCustomizationRepository(),
		}
	})
	return customizationInstance
}

func NewCustomizationService() *CustomizationService {
	return GetCustomizationService()
}

func (s *CustomizationService) GetBySubscriptionID(subscriptionID uint) (*model.SubscriptionCustomization, error) {
	c, err := s.repo.FindBySubscriptionID(subscriptionID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CustomizationService) Save(c *model.SubscriptionCustomization) error {
	return s.repo.Upsert(c)
}

func (s *CustomizationService) DeleteBySubscriptionID(subscriptionID uint) error {
	return s.repo.DeleteBySubscriptionID(subscriptionID)
}
