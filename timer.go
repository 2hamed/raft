package main

import (
	"math/rand"
	"time"
)

type RaftTimer struct {
	timeout       time.Duration
	timeoutSignal chan time.Time

	ticker *time.Ticker
}

func NewTimer() *RaftTimer {
	return &RaftTimer{
		timeout:       time.Duration(rand.Float64()*1000) + 300,
		timeoutSignal: make(chan time.Time),
	}
}

func (t *RaftTimer) TimeoutSignal() <-chan time.Time {
	return t.timeoutSignal
}

func (t *RaftTimer) StartTimer() {
	t.ticker = time.NewTicker(t.timeout * time.Millisecond)
	go func() {
		for {
			t.timeoutSignal <- <-t.ticker.C
		}
	}()
}

func (t *RaftTimer) Reset() {
	t.ticker.Reset(t.timeout)
}
