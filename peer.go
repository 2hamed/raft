package main

import "fmt"

type Peer struct {
	Host string
	Port int
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
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
