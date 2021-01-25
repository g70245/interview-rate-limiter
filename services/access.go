package services

import (
	"app/entities"
	"sync"
)

type Access struct {
	limiter *entities.RateLimiter
}

var singleton *Access
var once sync.Once

func GetAccessInstance() *Access {
	once.Do(func() {
		limiter := entities.CreateRateLimiter()
		singleton = &Access{limiter: &limiter}
	})
	return singleton
}

func (access *Access) Get(ip string, timestamp int64) (count int) {
	if tracker, ok := access.limiter.GetIPTracker(ip); ok {
		count = tracker.AddTrack(timestamp)
		return
	}

	tracker := entities.CreateIPTracker(ip, timestamp)
	if ok := access.limiter.AddIPTracker(&tracker); ok {
		count = 1
	}

	return
}

func (access *Access) Prune() {
	access.limiter.Prune()
}
