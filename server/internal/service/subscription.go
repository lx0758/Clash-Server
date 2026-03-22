package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"clash-server/internal/config"
	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/pkg/userinfo"

	"gopkg.in/yaml.v3"
)

type SubscriptionService struct {
	repo       *repository.SubscriptionRepository
	ruleRepo   *repository.RuleRepository
	scriptRepo *repository.ScriptRepository
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		repo:       repository.NewSubscriptionRepository(),
		ruleRepo:   repository.NewRuleRepository(),
		scriptRepo: repository.NewScriptRepository(),
	}
}

func (s *SubscriptionService) List() ([]model.Subscription, error) {
	return s.repo.FindAll()
}

type SubscriptionWithCounts struct {
	Subscription *model.Subscription `json:"subscription"`
	RuleCount    int                 `json:"rule_count"`
	ScriptCount  int                 `json:"script_count"`
}

func (s *SubscriptionService) ListWithCounts() ([]SubscriptionWithCounts, error) {
	subs, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	results := make([]SubscriptionWithCounts, len(subs))
	for i, sub := range subs {
		ruleCount, _ := s.ruleRepo.CountBySubscriptionID(sub.ID)
		scriptCount, _ := s.scriptRepo.CountBySubscriptionID(sub.ID)
		results[i] = SubscriptionWithCounts{
			Subscription: &sub,
			RuleCount:    ruleCount,
			ScriptCount:  scriptCount,
		}
	}
	return results, nil
}

func (s *SubscriptionService) Get(id uint) (*model.Subscription, error) {
	return s.repo.FindByID(id)
}

func (s *SubscriptionService) GetWithRelations(id uint) (*model.Subscription, []model.Rule, []model.Script, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	rules, err := s.ruleRepo.FindBySubscriptionID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	scripts, err := s.scriptRepo.FindBySubscriptionID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	return sub, rules, scripts, nil
}

func (s *SubscriptionService) Create(sub *model.Subscription) error {
	if sub.SourceType == "" {
		sub.SourceType = model.SourceTypeRemote
	}
	if sub.SourceType == model.SourceTypeRemote && sub.URL == "" {
		return errors.New("url is required for remote subscription")
	}
	if sub.SourceType == model.SourceTypeLocal {
		sub.URL = ""
		sub.Interval = 0
		sub.UseProxy = false
		sub.SkipCert = false
	} 
	if sub.SourceType == model.SourceTypeRemote {
		err := s.fetch(sub)
		if err != nil {
			return err
		}
	}
	count, err := s.repo.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		sub.Active = true
	}
	return s.repo.Create(sub)
}

func (s *SubscriptionService) Activate(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SetActive(id)
}

func (s *SubscriptionService) GetActive() (*model.Subscription, error) {
	sub, err := s.repo.GetActive()
	if err == nil {
		return sub, nil
	}
	subs, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return nil, nil
	}
	if err := s.repo.SetActive(subs[0].ID); err != nil {
		return nil, err
	}
	return &subs[0], nil
}

func (s *SubscriptionService) Update(sub *model.Subscription) error {
	if sub.SourceType == model.SourceTypeRemote {
		if sub.URL == "" {
			return errors.New("url is required for remote subscription")
		}
		err := s.fetch(sub)
		if err != nil {
			return err
		}
	}
	return s.repo.Update(sub)
}

func (s *SubscriptionService) UpdateContent(id uint, content string) error {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	sub.Content = content
	return s.repo.Update(sub)
}

func (s *SubscriptionService) Delete(id uint) error {
	if err := s.ruleRepo.DeleteBySubscriptionID(id); err != nil {
		return err
	}
	if err := s.scriptRepo.DeleteBySubscriptionID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

type RefreshResult struct {
	Subscription *model.Subscription
	Error        string
}

func (s *SubscriptionService) Refresh(id uint) (*RefreshResult, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sub.SourceType == model.SourceTypeLocal {
		return nil, errors.New("cannot refresh local subscription")
	}
	err = s.fetch(sub)
	if err != nil {
		return &RefreshResult{Subscription: sub, Error: err.Error()}, nil
	}
	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}
	return &RefreshResult{Subscription: sub}, nil
}

func (s *SubscriptionService) fetch(sub *model.Subscription) error {
	u, err := url.Parse(sub.URL)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("flag", "clash+verge+meta")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	ua := "clash-server clash-meta clash-verge"
	if sub.UserAgent != "" {
		ua = sub.UserAgent
	}
	req.Header.Set("User-Agent", ua)

	transport := &http.Transport{}
	if sub.UseProxy {
		u, err := url.Parse(fmt.Sprintf("http://localhost:%d", config.GetCoreConfig().MixedPort))
		if err == nil {
			transport.Proxy = http.ProxyURL(u)
		}
	}
	if sub.SkipCert {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch subscription")
	}

	userInfoHeader := resp.Header.Get("subscription-userinfo")
	if userInfoHeader != "" {
		info := userinfo.ParseSubscriptionUserinfo(userInfoHeader)
		if info != nil {
			sub.UploadUsed = info.UploadUsed
			sub.DownloadUsed = info.DownloadUsed
			sub.TotalTransfer = info.TotalTransfer
			sub.ExpireAt = info.ExpireAt
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	content := string(body)

	nodeCount := s.countNodes(content)
	if nodeCount == 0 {
		return errors.New("not a valid subscription")
	}
	sub.Content = content
	sub.NodeCount = nodeCount
	now := time.Now()
	sub.LastRefresh = &now
	return nil
}

func (s *SubscriptionService) countNodes(content string) int {
	var config map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return 0
	}
	proxies, ok := config["proxies"].([]interface{})
	if !ok {
		return 0
	}
	return len(proxies)
}

func (s *SubscriptionService) Parse(content string) (map[string]interface{}, error) {
	var config map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (s *SubscriptionService) GetMergedConfig(id uint) (map[string]interface{}, string, error) {
	merger := GetMergerService()
	config, err := merger.MergeForSubscription(id)
	if err != nil {
		return nil, "", err
	}
	yamlStr, err := merger.GenerateYAML(config)
	if err != nil {
		return nil, "", err
	}
	return config, yamlStr, nil
}
