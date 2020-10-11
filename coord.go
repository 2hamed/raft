package main

import (
	"encoding/json"
	"fmt"
)

type Coordinator struct {
	Self  Peer
	Peers Peers
}

func NewCoordinator(host string, port int) *Coordinator {
	return &Coordinator{
		Peers: make(Peers, 0),
		Self:  Peer{host, port},
	}
}

func (c *Coordinator) ProcessMessage(msg Message) error {

	switch msg.Ops {
	case "register", "reanounce", "propogate":
		c.peerInfoReceived(msg)
		break
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
		c.broadcastMessage(NewPropogateMessage(p.Host, p.Port))
		c.registerPeer(p)
		c.reanounceSelf()
		break
	case "reanounce":
		c.registerPeer(p)
		break
	case "propogate":
		c.registerPeer(p)
		c.reanounceSelf()
		break
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
	return peer.SendMessage(NewRegisterMessage("127.0.0.1", port))
}

func (c *Coordinator) reanounceSelf() error {
	return c.broadcastMessage(NewReanounceMessage(c.Self.Host, c.Self.Port))
}
