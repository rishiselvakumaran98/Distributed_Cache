package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"flag"

	cache "github.com/rishiselvakumaran98/Distributed_Cache/cache"
)
func main() {
	var (
		listenAddr = flag.String("listenAddr", ":3000", "listen address of the serv")
		leaderAddr = flag.String("leaderAddr", "", "listen address of the leader")
	)
	
	flag.Parse()
	opts := ServerOpts {
		ListenAddr: *listenAddr,
		IsLeader: true,
		LeaderAddr: *leaderAddr, // we need runFollower to add the argument for this 
	}
	// We make a simple client
	go func() {
		conn, err := net.Dial("tcp", *listenAddr) // activate a listener port in port 3000
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo Bar 25000000000")) // Once the port is reachable then write to the program that you reached port 3000
		time.Sleep(time.Second * 2)
		conn.Write([]byte("GET Foo"))
		buf := make([]byte, 1000)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	} () // we need to include parenthesis at end for performing go routine on a function

	server := NewServer(opts, cache.New())
	server.Start()
}