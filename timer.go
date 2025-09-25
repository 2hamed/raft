// RaftTimer implements a simple election timeout timer used by the coordinator
// to trigger leader elections when no heartbeat is received within a given
// duration.
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// RaftTimer implements a simple election timeout timer used by the coordinator
// to trigger leader elections when no heartbeat is received within a given
// duration.
type RaftTimer struct {
    timeout       time.Duration
    timeoutSignal chan time.Time

    ticker *time.Ticker
}
// NewTimer constructs a RaftTimer with a randomized timeout between 150–300ms.
// The jitter helps avoid simultaneous elections across multiple nodes.
func NewTimer() *RaftTimer {
    s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    return &RaftTimer{
        timeout:       time.Duration(r.Float64()*150+150) * time.Millisecond,
        timeoutSignal: make(chan time.Time),
    }
}

// TimeoutSignal exposes a read-only channel that fires when the election timer
// elapses.
func (t *RaftTimer) TimeoutSignal() <-chan time.Time {
    return t.timeoutSignal
}

// StartTimer starts the internal ticker and relays its ticks onto the
// timeoutSignal channel.
func (t *RaftTimer) StartTimer() {
    fmt.Println("Election timeout is", t.timeout)
    t.ticker = time.NewTicker(t.timeout)
    go func() {
        for {
            t.timeoutSignal <- <-t.ticker.C
        }
    }()
}

// Reset restarts the timer using the same timeout duration.
func (t *RaftTimer) Reset() {
    t.ticker.Reset(t.timeout)
}
