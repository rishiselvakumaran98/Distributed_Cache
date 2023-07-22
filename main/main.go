package main

import cache "github.com/rishiselvakumaran98/Distributed_Cache"

func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}
	server := NewServer(opts, cache.New())
	server.Start()
}