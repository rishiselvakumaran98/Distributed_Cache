package main

import (
	cache "github.com/rishiselvakumaran98/Distributed_Cache/cache"
	"fmt"
	"log"
)

// We can use a Raft Algorithm to do the selection of which
// cache node can be the leader to do write operations
type serverOpts struct {
	ListenAddr string
	isLeader bool
}

type Server struct {
	ServerOpts

	cache cache.CacheSet // we reference of 
}

func NewServer(opt ServerOpts, c cache.CacheSet) *Server {
	return &Server {
		ServerOpts: opts,
		cache: c,
	}
}

func (s *Server) Start() error {
	ln, err := net.listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("Listen Error: %s", err)
	}

	log.Printf("Server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

// Go routines
func (s *Server) handleConn(conn net.Conn) {
	// In this code, `make` is used to create a new slice with a specified length and capacity.
	buf := make([]byte, 2048)
	
	for {
		n, err := conn.Read(buf)
		if err != nil{
			log.Printf("Conn Read error: %s", err)
			break
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}