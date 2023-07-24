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
	LeaderAddr string
}

type Server struct {
	ServerOpts

	followers map[net.Conn] struct{

	}
	cache cache.CacheSet // we reference of
}

func NewServer(opts ServerOpts, c cache.CacheSet) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		// TODO: only allocate this when we are the leader
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("Listen Error: %s", err)
	}

	log.Printf("Server starting on port [%s]\n", s.ListenAddr)
	// if the current node is not a leader than it Dials to the leader to open new connection 
	if !s.IsLeader {
		go func(){
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Println("connected with Leader: ", s.LeaderAddr)
			if err != nil {
				log.Fatal(err)
			}
			s.handleConn(conn)
		}()
		
	}

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
	defer conn.Close()
	
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
		conn.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("Received Message: %s\n", msg.Cmd)
	var value []byte
	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCmd(conn, msg)
	case CMDGet:
		value, err = s.handleGetCmd(conn, msg)
		fmt.Printf("Get value: %s", string(value))
	}

	if err != nil {
		fmt.Printf("Failed to handle %s command with error: %s \n", msg.Cmd, err)
		conn.Write([]byte(err.Error()))
		return
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

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) ([]byte, error) {
	// if the server cannot get the key from the cache
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	// must be mutex protected
	for conn := range s.followers {
		_, err := conn.Write(msg.ToBytes())
		if err != nil {
			log.Printf("Write to followers failed %s", err)
			continue
		}
	}
	return nil
}
