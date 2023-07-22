package main

import cache "distributed_cache"

func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}
	server := NewServer(opts, cache.New())
	server.Start()
}