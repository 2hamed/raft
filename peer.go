package raft

import (
	"fmt"
	"net"
)

type peer struct {
	Host string
	Port int
}

func (p *peer) String() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func (p *peer) Equals(other peer) bool {
	return p.Host == other.Host && p.Port == other.Port
}

func (p *peer) SendMessage(msg message) error {
	raddr, err := net.ResolveUDPAddr("udp", p.String())
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(msg.Json())

	return err
}

type Peers []peer

func (peers *Peers) Contains(peer peer) bool {
	for _, p := range *peers {
		if p.Host == peer.Host && p.Port == peer.Port {
			return true
		}
	}
	return false
}

func (peers *Peers) PrintInfo() {
	fmt.Println("Connected Peers:")
	for _, p := range *peers {
		fmt.Printf("%s: %d\n", p.Host, p.Port)
	}
}

func (peers *Peers) BroadcastMessage(msg message) (err error) {
	failedPeers := make([]int, 0)
	for i, p := range *peers {
		err = p.SendMessage(msg)
		if err != nil {
			fmt.Printf("Sending message to peer %s failed: %v\n", p.String(), err)
			failedPeers = append(failedPeers, i)
		}
	}
	return err
}

func (peers *Peers) Quorum() int {
	return len(*peers)/2 + 1
}
