package entities

import (
	"testing"
	"time"
)

func TestAddWithOnlyOneTrack(t *testing.T) {
	tracker := IPTracker{}

	count := tracker.AddTrack(100001)
	if count != 1 || tracker.timestamps[0] != 100001 {
		t.Fail()
	}
}

func TestMaxCapacityOfTracks(t *testing.T) {
	tracker := IPTracker{}

	var count int
	for i := 0; i < 100; i++ {
		count = tracker.AddTrack(100001)
	}
	if count > MAX_TRACKS {
		t.Fail()
	}
}

func TestAddWithTracksExceedingMaxInSameWindow(t *testing.T) {
	tracker := IPTracker{}

	for i := 0; i < 60; i++ {
		tracker.AddTrack(time.Now().Unix())
	}
	last := time.Now().Unix()
	tracker.AddTrack(last)
	count := tracker.AddTrack(time.Now().Unix())

	if count > MAX_TRACKS {
		t.Error("Should reach max capacity")
	}

	if tracker.timestamps[MAX_TRACKS-1] != last {
		t.Error("Incorrect last track")
	}
}

func TestTrackWindowShifting(t *testing.T) {
	tracker := IPTracker{}
	givenTracksNum := 40

	tracker.AddTrack(100001)
	tracker.AddTrack(100001)
	for i := 0; i < givenTracksNum-1; i++ {
		tracker.AddTrack(time.Now().Unix())
	}
	last := time.Now().Unix()
	count := tracker.AddTrack(last)

	if count != givenTracksNum {
		t.Error("Incorrect length of track records")
	}

	if tracker.timestamps[givenTracksNum-1] != last {
		t.Error("Incorrect last track")
	}
}

func TestCreate(t *testing.T) {
	ip := "127.0.0.1"
	ts := time.Now().Unix()
	tracker := CreateIPTracker(ip, ts)

	if tracker.ip != ip {
		t.Error("Incorrect IP")
	}

	if tracker.timestamps[0] != ts {
		t.Error("Incorrect timestamp")
	}
}

func TestGetLastTimestamp(t *testing.T) {
	tracker := IPTracker{}

	tracker.AddTrack(time.Now().Unix())
	tracker.AddTrack(time.Now().Unix())
	tracker.AddTrack(time.Now().Unix())
	last := time.Now().Unix()
	tracker.AddTrack(last)

	if tracker.GetLastTrack() != last {
		t.Fail()
	}
}
