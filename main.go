package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
)

func main() {
	peerHost := flag.String("phost", "127.0.0.1", "peer host")
	peerPort := flag.Int("pport", 0, "peer port")
	port := flag.Int("port", 3000, "port to listen on")

	flag.Parse()

	coord := NewCoordinator("127.0.0.1", *port)

	forever := make(chan bool)

	startUDPServer(*port, coord)

	if *peerPort != 0 {
		coord.joinCluster(*port, Peer{*peerHost, *peerPort})
	}

	<-forever
}

func startUDPServer(port int, coord *Coordinator) {
	fmt.Println("UDP server listening on port: ", port)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
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

			coord.ProcessMessage(message)
		}
	}()
}
