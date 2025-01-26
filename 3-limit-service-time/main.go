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
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) AddTime(s int64) int64 {
	return atomic.AddInt64(&u.TimeUsed, s)
}

func (u *User) GetTimeUsed() int64 {
	return atomic.LoadInt64(&u.TimeUsed)
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	processChan := make(chan struct{})

	go func() {
		process()
		processChan <- struct{}{}
	}()

	if u.IsPremium {
		<-processChan
		return true
	}

	if u.GetTimeUsed() >= 10 {
		return false
	}

	for {
		select {
		case <-time.Tick(time.Second):
			if u.AddTime(1) >= 10 {
				return false
			}
		case <-processChan:
			if u.GetTimeUsed() >= 10 {
				return false
			}
			return true
		}
	}
}

func main() {
	RunMockServer()
}
