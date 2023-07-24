package main

import (
	"fmt"
	"log"
	"net"
	"time"

	cache "github.com/rishiselvakumaran98/Distributed_Cache/cache"
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
		conn.Write([]byte("SET Foo Bar 2500000")) // Once the port is reachable then write to the program that you reached port 3000
		time.Sleep(time.Second * 2)
		conn.Write([]byte("GET Foo"))
		buf := make([]byte, 1000)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	} () // we need to include parenthesis at end for performing go routine on a function

	server := NewServer(opts, cache.New())
	server.Start()
}