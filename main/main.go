package main

import cache "rishiselvakumaran98.com/cache"

func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}
	server := NewServer(opts, cache.New())
	server.Start()
}