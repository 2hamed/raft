package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type RaftServer struct {
	options Options

	coord *Coordinator
}

var started = false

func NewRaftServer(options ...OptionFunc) *RaftServer {
	opts := &Options{}
	for _, optFunc := range options {
		optFunc(opts)
	}

	return &RaftServer{
		options: *opts,
		coord:   NewCoordinator(opts.listenAddr, opts.listenPort),
	}
}

func (r *RaftServer) Start(ctx context.Context) error {
	fmt.Printf("UDP server listening on %s:%d .", r.options.listenAddr, r.options.listenPort)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(r.options.listenAddr), Port: r.options.listenPort})
	if err != nil {
		return fmt.Errorf("Failed starting the UDP server, %w", err)
	}

	msgChan := make(chan []byte)

	go func() {
		for {
			msg := make([]byte, 1024)
			n, _, err := listener.ReadFromUDP(msg)
			if err != nil {
				fmt.Println("Error reading message from UDP:", err)
				continue
			}
			msgChan <- msg[:n]
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down raft...")
			return nil
		case m := <-msgChan:
			var message Message
			if err := json.Unmarshal(m, &message); err != nil {
				fmt.Printf("Error in unmarshaling the message: %v\n", err)
				continue
			}

			r.coord.ProcessMessage(message)
		}
	}
}

func (r *RaftServer) JoinCluster(ctx context.Context, peerHost string, peerPort int) error {
	if !started {
		return errors.New("Call Start first")
	}

	return r.coord.joinCluster(r.options.listenAddr, r.options.listenPort, Peer{peerHost, peerPort})
}
