package services

import (
	"strconv"
	"testing"
	"time"
)

func TestGetInstance(t *testing.T) {
	access1 := GetAccessInstance()
	access2 := GetAccessInstance()

	ch := make(chan *Access)
	go func(ch chan *Access) {
		ch <- GetAccessInstance()
	}(ch)
	access3 := <-ch

	if access1 != access2 || access2 != access3 || access1 != access3 {
		t.Fail()
	}
}

func TestPrune(t *testing.T) {
	access := GetAccessInstance()

	ip := "127.0.0.1"
	timestamp := time.Now().Unix() - 59
	access.Get(ip, timestamp)
	access.Get(ip, timestamp)

	time.Sleep(2000 * time.Millisecond)
	access.Prune()

	if size := access.limiter.Size(); size != 0 {
		t.Fail()
	}
}

func TestGetWithOnlyOneRequest(t *testing.T) {
	access := GetAccessInstance()

	ip := "127.0.0.1"
	timestamp := time.Now().Unix()
	if count := access.Get(ip, timestamp); count != 1 {
		t.Fail()
	}
}

func TestGetWithRequestsMoreThan60InOneMinute(t *testing.T) {
	access := GetAccessInstance()

	ip := "127.0.0.1"
	timestamp := time.Now().Unix()

	var count int
	for i := 0; i < 80; i++ {
		count = access.Get(ip, timestamp)
	}
	if count != 61 {
		t.Fail()
	}
}

func TestGetWithMultiIPRequests(t *testing.T) {
	access := GetAccessInstance()

	ips := make([]string, 6)
	counts := make([]int, 6)
	for i := 1; i < 6; i++ {
		ips[i] = "127.0.0." + strconv.Itoa(i+1)
	}

	var count int
	for i, ip := range ips {
		for j := 0; j < i*10; j++ {
			count = access.Get(ip, time.Now().Unix())
		}
		counts[i] = count
	}

	for i, c := range counts {
		if c != i*10 {
			t.Fail()
		}
	}
}
