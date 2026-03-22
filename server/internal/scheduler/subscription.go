package scheduler

import (
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

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[sub.ID]; exists {
		return
	}

	interval := time.Duration(sub.Interval) * time.Minute
	timer := time.AfterFunc(interval, func() {
		s.refreshSubscription(sub)
		s.scheduleRefresh(sub)
	})

	s.timers[sub.ID] = timer
}

func (s *SubscriptionScheduler) refreshSubscription(sub *model.Subscription) {
	select {
	case <-s.stopChan:
		return
	default:
	}

	if _, err := s.subService.Refresh(sub.ID); err != nil {
		return
	}

	if sub.Active {
		s.coreService.ApplyConfig()
	}
}
