package main

import (
	"fmt"
	"net"
)

// Peer represents a single cluster node identified by host and port.
type Peer struct {
	Host string
	Port int
}

// String returns the host:port representation of the peer.
func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

// Equals compares two peers for equality based on host and port.
func (p *Peer) Equals(other Peer) bool {
	return p.Host == other.Host && p.Port == other.Port
}

// SendMessage sends a UDP message to this peer.
func (p *Peer) SendMessage(msg Message) error {
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

// Peers is a convenience alias for a slice of Peer.
type Peers []Peer

// Contains reports whether the given peer is already present in the set.
func (peers *Peers) Contains(peer Peer) bool {
	for _, p := range *peers {
		if p.Host == peer.Host && p.Port == peer.Port {
			return true
		}
	}
	return false
}

// PrintInfo logs the currently connected peers to stdout.
func (peers *Peers) PrintInfo() {
	fmt.Println("Connected Peers:")
	for _, p := range *peers {
		fmt.Printf("%s: %d\n", p.Host, p.Port)
	}
}

// BroadcastMessage sends the provided message to all peers, returning the last
// error encountered (if any). Errors are logged and not retried.
func (peers *Peers) BroadcastMessage(msg Message) (err error) {
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

// Quorum returns the number of votes required to reach a majority, assuming the
// local node plus all known peers comprise the cluster.
func (peers *Peers) Quorum() int {
	return (len(*peers)+1)/2 + 1
}
