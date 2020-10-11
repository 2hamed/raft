package main

import "encoding/json"

type Coordinator struct {
	Peers Peers
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		Peers: make(Peers, 0),
	}
}

func (c *Coordinator) ProcessMessage(msg Message) error {

	switch msg.Ops {
	case "register":
		var p Peer
		if err := json.Unmarshal([]byte(msg.Payload), &p); err != nil {
			return err
		}

		if err := c.registerPeer(p); err != nil {
			return err
		}

		break
	}

	return nil
}

func (c *Coordinator) registerPeer(peer Peer) error {
	if !c.Peers.Contains(peer) {
		c.Peers = append(c.Peers, peer)
	}
	c.Peers.PrintInfo()
	return nil
}
