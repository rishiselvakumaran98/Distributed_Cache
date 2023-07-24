package main

import (
	"context"
	"fmt"
	"log"
	"net"

	cache "github.com/rishiselvakumaran98/Distributed_Cache/cache"
)

// We can use a Raft Algorithm to do the selection of which
// cache node can be the leader to do write operations
type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts

	cache cache.CacheSet // we reference of
}

func NewServer(opts ServerOpts, c cache.CacheSet) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
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
	// defer the connection from being closed when the listening port is still up
	defer func() {
		conn.Close()
	}()
	// In this code, `make` is used to create a new slice with a specified length and capacity.
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Conn Read error: %s", err)
			break
		}
		msg := buf[:n]
		go s.handleCommand(conn, msg)
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("Failed to parse the command", err)
		// Respond
		return
	}
	switch msg.Cmd {
	case CMDSet:
		if err := s.handleSetCmd(conn, msg); err != nil {
			// respond
			return
		}
	}

}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	// if the cache cannot set the key and value bytes then return the err
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)
	
	fmt.Println("Handling set command", msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
