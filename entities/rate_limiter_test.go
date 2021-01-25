package entities

import (
	"testing"
	"time"
)

func TestAddTracker(t *testing.T) {
	ip := "127.0.0.1"
	timestamp := time.Now().Unix()
	tracker := CreateIPTracker(ip, timestamp)

	limiter := CreateRateLimiter()
	if ok := limiter.AddIPTracker(&tracker); ok {
		if len(limiter.trackers) != 1 || limiter.trackers[ip].ip != ip || limiter.trackers[ip].timestamps[0] != timestamp {
			t.Error("Incorrect attributes")
		}
		if limiter.trackers[ip] != &tracker {
			t.Error("Should be the same object")
		}
	} else {
		t.Error("AddTracker Failed")
	}
}

func TestGetTracker(t *testing.T) {
	ip := "127.0.0.1"
	timestamp := time.Now().Unix()
	tracker := CreateIPTracker(ip, timestamp)

	limiter := CreateRateLimiter()
	limiter.AddIPTracker(&tracker)
	if fetchedTracker, _ := limiter.GetIPTracker(ip); fetchedTracker != &tracker {
		t.Error("Should be the same object")
	}
}

func TestPrune(t *testing.T) {
	ip1 := "127.0.0.1"
	timestamp1 := time.Now().Unix() - (SECONDS_PER_TRACK_WINDOW + 1)
	tracker1 := CreateIPTracker(ip1, timestamp1)

	ip2 := "127.0.0.2"
	timestamp2 := time.Now().Unix()
	tracker2 := CreateIPTracker(ip2, timestamp2)

	limiter := CreateRateLimiter()
	limiter.AddIPTracker(&tracker1)
	limiter.AddIPTracker(&tracker2)

	limiter.Prune()
	if len(limiter.trackers) != 1 {
		t.Error("Prune failed")
	}
	if tracker, ok := limiter.GetIPTracker(ip2); !ok {
		t.Error("Should delete the ip tracker that has not been active in 60 seconds")
	} else if tracker != &tracker2 {
		t.Error("Should be the same object")
	}
}

// AddTracker no longer trigers prune
/* func TestAddPrune(t *testing.T) {
	limiter := CreateRateLimiter()

	ip := "127.0.0.1"
	timestamp := time.Now().Unix() - (SECONDS_PER_TRACK_WINDOW + 1)
	prunedTracker := CreateIPTracker(ip, timestamp)
	limiter.AddIPTracker(&prunedTracker)

	for i := 2; i < 6010; i++ {
		tracker := CreateIPTracker("127.0.0."+strconv.Itoa(i), time.Now().Unix())
		limiter.AddIPTracker(&tracker)
	}

	if len(limiter.trackers) != 6000 {
		t.Error("Prune failed")
	}

	if _, ok := limiter.GetIPTracker(ip); ok {
		t.Error("Should delete the ip tracker that has not been active in 60 seconds")
	}
} */
