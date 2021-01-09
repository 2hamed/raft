package raft

import (
	"flag"
)

func main() {
	peerHost := flag.String("phost", "127.0.0.1", "peer host")
	peerPort := flag.Int("pport", 0, "peer port")
	addr := flag.String("addr", "127.0.0.1", "network address to listen on")
	port := flag.Int("port", 3000, "port to listen on")

	flag.Parse()

	coord := NewCoordinator(*addr, *port)

	forever := make(chan bool)

	startUDPServer(*addr, *port, coord)

	if *peerPort != 0 {
		coord.joinCluster(*port, Peer{*peerHost, *peerPort})
	}

	<-forever
}

func startUDPServer(addr string, port int, coord *Coordinator) {
	
}
