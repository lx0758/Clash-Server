package scheduler

import (
	"log"
	"sync"
	"time"

	"clash-server/internal/model"
	"clash-server/internal/repository"
	"clash-server/internal/service"
)

type SubscriptionScheduler struct {
	subService  *service.SubscriptionService
	coreService *service.CoreService
	repo        *repository.SubscriptionRepository
	timers      map[uint]*time.Timer
	mu          sync.RWMutex
	stopChan    chan struct{}
}

func NewSubscriptionScheduler() *SubscriptionScheduler {
	return &SubscriptionScheduler{
		subService:  service.NewSubscriptionService(),
		coreService: service.GetCoreService(),
		repo:        repository.NewSubscriptionRepository(),
		timers:      make(map[uint]*time.Timer),
		stopChan:    make(chan struct{}),
	}
}

func (s *SubscriptionScheduler) Start() error {
	s.mu.Lock()
	s.stopChan = make(chan struct{})
	s.timers = make(map[uint]*time.Timer)
	s.mu.Unlock()

	subs, err := s.repo.FindAll()
	if err != nil {
		return err
	}
	for i := range subs {
		s.scheduleRefresh(&subs[i])
	}
	return nil
}

func (s *SubscriptionScheduler) Stop() {
	close(s.stopChan)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, timer := range s.timers {
		timer.Stop()
	}
	s.timers = make(map[uint]*time.Timer)
}

func (s *SubscriptionScheduler) AddSubscription(sub *model.Subscription) {
	s.scheduleRefresh(sub)
}

func (s *SubscriptionScheduler) UpdateSubscription(sub *model.Subscription) {
	s.RemoveSubscription(sub.ID)
	s.scheduleRefresh(sub)
}

func (s *SubscriptionScheduler) RemoveSubscription(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if timer, ok := s.timers[id]; ok {
		timer.Stop()
		delete(s.timers, id)
	}
}

func (s *SubscriptionScheduler) scheduleRefresh(sub *model.Subscription) {
	if sub.Interval <= 0 {
		return
	}
	if sub.SourceType != model.SourceTypeRemote {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[sub.ID]; exists {
		return
	}

	var nextRefresh time.Duration
	now := time.Now()

	if sub.LastRefresh != nil {
		nextRefreshTime := sub.LastRefresh.Add(time.Duration(sub.Interval) * time.Minute)
		if now.After(nextRefreshTime) || now.Equal(nextRefreshTime) {
			nextRefresh = time.Duration(5 * time.Minute)
		} else {
			nextRefresh = time.Until(nextRefreshTime)
		}
	} else {
		nextRefresh = time.Duration(sub.Interval) * time.Minute
	}

	subID := sub.ID
	timer := time.AfterFunc(nextRefresh, func() {
		updatedSub, err := s.repo.FindByID(subID)
		if err != nil {
			log.Printf("[Scheduler] Failed to find subscription %d: %v", subID, err)
			return
		}
		s.mu.Lock()
		delete(s.timers, subID)
		s.mu.Unlock()
		s.refreshSubscription(updatedSub)
		s.scheduleRefresh(updatedSub)
	})

	s.timers[sub.ID] = timer
}

func (s *SubscriptionScheduler) refreshSubscription(sub *model.Subscription) {
	select {
	case <-s.stopChan:
		return
	default:
	}

	result, err := s.subService.Refresh(sub.ID)
	if err != nil {
		log.Printf("[Scheduler] Failed to refresh subscription %d: %v", sub.ID, err)
		return
	}
	if result.Error != "" {
		log.Printf("[Scheduler] Subscription %d refresh error: %s", sub.ID, result.Error)
	}

	if sub.Active {
		s.coreService.ApplyConfig()
	}
}
