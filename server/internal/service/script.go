package service

import (
	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/pkg/script"
)

type ScriptService struct {
	repo   *repository.ScriptRepository
	engine *script.Engine
}

func NewScriptService() *ScriptService {
	return &ScriptService{
		repo:   repository.NewScriptRepository(),
		engine: script.NewEngine(),
	}
}

func (s *ScriptService) List(subscriptionID uint) ([]model.Script, error) {
	return s.repo.FindBySubscriptionID(subscriptionID)
}

func (s *ScriptService) Get(id uint) (*model.Script, error) {
	return s.repo.FindByID(id)
}

func (s *ScriptService) GetBySubscription(id, subscriptionID uint) (*model.Script, error) {
	return s.repo.FindByIDAndSubscriptionID(id, subscriptionID)
}

func (s *ScriptService) Create(scriptModel *model.Script) error {
	return s.repo.Create(scriptModel)
}

func (s *ScriptService) Update(scriptModel *model.Script) error {
	return s.repo.Update(scriptModel)
}

func (s *ScriptService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ScriptService) DeleteBySubscription(subscriptionID uint) error {
	return s.repo.DeleteBySubscriptionID(subscriptionID)
}

func (s *ScriptService) Test(id uint, config map[string]interface{}) (map[string]interface{}, error) {
	scriptModel, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.engine.Execute(scriptModel.Content, config)
}

func (s *ScriptService) TestBySubscription(id, subscriptionID uint, config map[string]interface{}) (map[string]interface{}, error) {
	scriptModel, err := s.repo.FindByIDAndSubscriptionID(id, subscriptionID)
	if err != nil {
		return nil, err
	}
	return s.engine.Execute(scriptModel.Content, config)
}

func (s *ScriptService) ExecuteEnabled(subscriptionID uint, config map[string]interface{}) (map[string]interface{}, error) {
	scripts, err := s.repo.FindEnabledBySubscriptionID(subscriptionID)
	if err != nil {
		return nil, err
	}
	result := config
	for _, script := range scripts {
		result, err = s.engine.Execute(script.Content, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
