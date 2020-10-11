package main

import (
	"fmt"
	"math/rand"
	"time"
)

type RaftTimer struct {
	timeout       time.Duration
	timeoutSignal chan time.Time

	ticker *time.Ticker
}

func NewTimer() *RaftTimer {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return &RaftTimer{
		timeout:       time.Duration(r.Float64()*150+150) * time.Millisecond,
		timeoutSignal: make(chan time.Time),
	}
}

func (t *RaftTimer) TimeoutSignal() <-chan time.Time {
	return t.timeoutSignal
}

func (t *RaftTimer) StartTimer() {
	fmt.Println("Election timeout is", t.timeout)
	t.ticker = time.NewTicker(t.timeout)
	go func() {
		for {
			t.timeoutSignal <- <-t.ticker.C
		}
	}()
}

func (t *RaftTimer) Reset() {
	t.ticker.Reset(t.timeout)
}
