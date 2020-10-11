package main

import (
	"fmt"
	"net"
)

type Peer struct {
	Host string
	Port int
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func (p *Peer) Equals(other Peer) bool {
	return p.Host == other.Host && p.Port == other.Port
}

func (p *Peer) SendMessage(msg Message) error {
	raddr, err := net.ResolveUDPAddr("udp", p.String())
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	_, err = conn.Write(msg.Json())

	return nil
}

type Peers []Peer

func (peers *Peers) Contains(peer Peer) bool {
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

func (peers *Peers) BroadcastMessage(msg Message) (err error) {
	for _, p := range *peers {
		err = p.SendMessage(msg)
	}
	return err
}
