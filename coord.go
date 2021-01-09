package raft

import (
	"encoding/json"
	"fmt"
	"time"
)

type Coordinator struct {
	Self  Peer
	Peers Peers
	Timer *RaftTimer

	heartBeat   chan struct{}
	isElected   bool
	isCandidate bool

	votes int
}

func NewCoordinator(host string, port int) *Coordinator {
	c := &Coordinator{
		Peers:     make(Peers, 0),
		Self:      Peer{host, port},
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
				c.PromoteSelf()
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

func (c *Coordinator) ProcessMessage(msg Message) error {

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

func (c *Coordinator) broadcastMessage(msg Message) error {
	return c.Peers.BroadcastMessage(msg)
}

func (c *Coordinator) peerInfoReceived(msg Message) error {
	var p Peer
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
func (c *Coordinator) registerPeer(peer Peer) error {
	if !c.Peers.Contains(peer) {
		c.Peers = append(c.Peers, peer)
	}

	return nil
}

func (c *Coordinator) joinCluster(port int, peer Peer) error {
	fmt.Printf("Joining peer on %s:%d\n", peer.Host, peer.Port)
	return peer.SendMessage(NewRegisterMessage("127.0.0.1", port).WithSender(c.Self))
}

func (c *Coordinator) reanounceSelf() error {
	return c.broadcastMessage(NewReanounceMessage(c.Self.Host, c.Self.Port).WithSender(c.Self))
}

func (c *Coordinator) PromoteSelf() error {
	fmt.Println("Starting election...")
	c.votes = 1
	return c.broadcastMessage(NewPromoteMessage().WithSender(c.Self))
}

func (c *Coordinator) sendHeartbeat() error {
	return c.broadcastMessage(NewHeartbeatMessage().WithSender(c.Self))
}

func (c *Coordinator) sendVote(target Peer) error {
	return target.SendMessage(NewVoteMessage().WithSender(c.Self))
}
