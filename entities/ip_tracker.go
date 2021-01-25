package entities

import (
	"sync"
	"time"
)

const SECONDS_PER_TRACK_WINDOW = 60
const MAX_TRACKS = 61

type IPTracker struct {
	ip         string
	mutex      sync.Mutex
	timestamps []int64
}

func CreateIPTracker(ip string, timestamp int64) IPTracker {
	tracker := IPTracker{ip: ip}
	tracker.AddTrack(timestamp)
	return tracker
}

func (tracker *IPTracker) AddTrack(timestamp int64) (count int) {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()

	shiftTrackWindow(tracker)

	count = len(tracker.timestamps)
	if count == MAX_TRACKS {
		return
	}
	tracker.timestamps = append(tracker.timestamps, timestamp)
	return count + 1
}

func (tracker *IPTracker) GetLastTrack() (timestamp int64) {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()

	return tracker.timestamps[len(tracker.timestamps)-1]
}

func shiftTrackWindow(tracker *IPTracker) {
	startTimestamp := time.Now().Unix() - SECONDS_PER_TRACK_WINDOW

	for i, v := range tracker.timestamps {
		if v >= startTimestamp {
			tracker.timestamps = tracker.timestamps[i:]
			break
		}
	}
}
