package raft

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type RaftServer struct {
	options Options

	coord *Coordinator
}

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

func (r *RaftServer) Start(ctx context.Context) {
	fmt.Printf("UDP server listening on %s:%d .", r.options.listenAddr, r.options.listenPort)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(r.options.listenAddr), Port: r.options.listenPort})
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down raft...")
			return
		default:
			var msg = make([]byte, 1024)
			n, addr, err := listener.ReadFromUDP(msg)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Read %d bytes from %v, Body: %s\n", n, addr, (msg))
			var message Message
			if err := json.Unmarshal(msg[:n], &message); err != nil {
				fmt.Printf("Error in unmarshaling the message: %v\n", err)
				continue
			}

			r.coord.ProcessMessage(message)
		}
	}
}
