package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	peerHost := flag.String("phost", "127.0.0.1", "peer host")
	peerPort := flag.Int("pport", 0, "peer port")
	port := flag.Int("port", 3000, "port to listen on")

	flag.Parse()

	forever := make(chan bool)

	startUDPServer(*port)

	fmt.Printf("Joining peer on %s:%d\n", *peerHost, *peerPort)

	<-forever
}

func startUDPServer(port int) {
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
			fmt.Printf("Read %d bytes from %v, Body: %s", n, addr, (msg))
			time.Sleep(1 * time.Second)
		}
	}()
}
