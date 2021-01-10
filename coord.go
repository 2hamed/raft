package raft

import (
	"encoding/json"
	"fmt"
	"time"
)

type coordinator struct {
	Self  peer
	Peers Peers
	Timer *RaftTimer

	heartBeat   chan struct{}
	isElected   bool
	isCandidate bool

	votes int
}

func NewCoordinator(host string, port int) *coordinator {
	c := &coordinator{
		Peers:     make(Peers, 0),
		Self:      peer{host, port},
		Timer:     NewTimer(),
		heartBeat: make(chan struct{}, 10),
	}

	go func() {
		for {
			select {
			case <-c.heartBeat:
				c.Timer.Reset()
				if c.isElected {
					c.sendHeartbeat()
				}
			case <-c.Timer.TimeoutSignal():
				c.isCandidate = true
				c.promoteSelf()
			}
		}
	}()

	go func() {
		ticker := time.Tick(50 * time.Millisecond)
		for {
			<-ticker
			if c.isElected {
				c.heartBeat <- struct{}{}
			}
		}
	}()

	c.Timer.StartTimer()

	return c
}

func (c *coordinator) ProcessMessage(msg message) error {

	switch msg.Ops {
	case "register", "reanounce", "propogate":
		c.peerInfoReceived(msg)
	case "beat":
		c.heartBeat <- struct{}{}
	case "promote":
		if !c.isCandidate && !c.isElected {
			c.sendVote(msg.Sender)
		}
	case "vote":
		c.votes++
		if c.votes >= c.Peers.Quorum() {
			fmt.Println("I am the elected leader.")
			c.isElected = true
		}
	}

	return nil
}

func (c *coordinator) broadcastMessage(msg message) error {
	return c.Peers.BroadcastMessage(msg)
}

func (c *coordinator) peerInfoReceived(msg message) error {
	var p peer
	if err := json.Unmarshal([]byte(msg.Payload), &p); err != nil {
		return err
	}

	if c.Self.Equals(p) {
		return nil
	}

	switch msg.Ops {
	case "register":
		c.broadcastMessage(NewPropogateMessage(p.Host, p.Port).WithSender(c.Self))
		c.registerPeer(p)
		c.reanounceSelf()
	case "reanounce":
		c.registerPeer(p)
	case "propogate":
		c.registerPeer(p)
		c.reanounceSelf()
	}
	c.Peers.PrintInfo()
	return nil
}
func (c *coordinator) registerPeer(peer peer) error {
	if !c.Peers.Contains(peer) {
		c.Peers = append(c.Peers, peer)
	}

	return nil
}

func (c *coordinator) joinCluster(nodeListenAddr string, nodePort int, peer peer) error {
	fmt.Printf("Joining peer on %s:%d\n", peer.Host, peer.Port)
	return peer.SendMessage(NewRegisterMessage(nodeListenAddr, nodePort).WithSender(c.Self))
}

func (c *coordinator) reanounceSelf() error {
	return c.broadcastMessage(NewReanounceMessage(c.Self.Host, c.Self.Port).WithSender(c.Self))
}

func (c *coordinator) promoteSelf() error {
	fmt.Println("Starting election...")
	c.votes = 1
	return c.broadcastMessage(NewPromoteMessage().WithSender(c.Self))
}

func (c *coordinator) sendHeartbeat() error {
	return c.broadcastMessage(NewHeartbeatMessage().WithSender(c.Self))
}

func (c *coordinator) sendVote(target peer) error {
	return target.SendMessage(NewVoteMessage().WithSender(c.Self))
}
