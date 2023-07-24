package main

import (
	cache "github.com/rishiselvakumaran98/Distributed_Cache/cache"
	"log"
	"net"
)
func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}
	// We make a simple client
	go func() {
		conn, err := net.Dial("tcp", ":3000") // activate a listener port in port 3000
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo Bar 2500")) // Once the port is reachable then write to the program that you reached port 3000
	} () // we need to include parenthesis at end for performing go routine on a function

	server := NewServer(opts, cache.New())
	server.Start()
}