//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in milliseconds
}

const (
	MaxFreeMs   = 10000
	CheckPeriod = 100
)

// cache saves consumed time for each user in microseconds
type TimeTracker struct {
	m     sync.Mutex
	users map[int]*User
}

var tracker = TimeTracker{
	m:     sync.Mutex{},
	users: make(map[int]*User),
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	// create new user if wasn't registered
	tracker.m.Lock()
	if user, ok := tracker.users[u.ID]; !ok {
		tracker.users[u.ID] = u
	} else {
		if user.TimeUsed >= MaxFreeMs {
			return false
		}
	}
	tracker.m.Unlock()

	ticker := time.NewTicker(time.Duration(CheckPeriod) * time.Millisecond)
	defer ticker.Stop()

	// start processing
	done := make(chan struct{}, 1)
	go func() {
		process()
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			return true
		case <-ticker.C:
			tracker.m.Lock()
			user := tracker.users[u.ID]
			user.TimeUsed += CheckPeriod
			if user.TimeUsed >= MaxFreeMs {
				tracker.m.Unlock()
				return false
			}
			tracker.m.Unlock()
		}
	}
}

func main() {
	go func() {
		seconds := 0
		for range time.Tick(time.Second) {
			seconds++
			fmt.Println("Time: ", seconds)
		}
	}()
	RunMockServer()
}
