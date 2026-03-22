package handler

import (
	"clash-server/internal/scheduler"
)

var subScheduler *scheduler.SubscriptionScheduler

func SetSubscriptionScheduler(s *scheduler.SubscriptionScheduler) {
	subScheduler = s
}

func GetSubscriptionScheduler() *scheduler.SubscriptionScheduler {
	return subScheduler
}
