package entities

import (
	"sync"
	"time"
)

const MAX_IP_ADDRESSES = 6000

type RateLimiter struct {
	mutex    sync.RWMutex
	trackers map[string]*IPTracker
}

func CreateRateLimiter() (limiter RateLimiter) {
	limiter = RateLimiter{}
	limiter.trackers = make(map[string]*IPTracker)
	return
}

func (r *RateLimiter) AddIPTracker(tracker *IPTracker) (ok bool) {
	if len(r.trackers) == MAX_IP_ADDRESSES {
		return
		/* if hasDecreased := r.Prune(); !hasDecreased {
			return
		} */
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.trackers[tracker.ip] = tracker
	ok = true
	return
}

func (r *RateLimiter) GetIPTracker(ip string) (tracker *IPTracker, ok bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	t, ok := r.trackers[ip]
	return t, ok
}

func (r *RateLimiter) Prune() (hasDecreased bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	lengthBeforePruning := len(r.trackers)
	currentTimestamp := time.Now().Unix()

	for k, v := range r.trackers {
		if v.GetLastTrack() < currentTimestamp-SECONDS_PER_TRACK_WINDOW {
			delete(r.trackers, k)
		}
	}

	if len(r.trackers) < lengthBeforePruning {
		hasDecreased = true
	}
	return
}

func (r *RateLimiter) Size() int {
	return len(r.trackers)
}
