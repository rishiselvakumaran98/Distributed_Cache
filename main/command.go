package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Object used for declaring commands for the Cache
type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
	CMDDelete Command = "DELETE"
)

type Message struct {
	Cmd Command
	Key []byte
	Value []byte
	TTL time.Duration
}

func (m *Message) ToBytes() []byte {
	switch m.Cmd {
	case CMDSet:
		cmd := fmt.Sprintf("%s %s %s %s", m.Cmd, m.Key, m.Value, m.TTL)
		return []byte(cmd)
	case CMDGet:
		cmd := fmt.Sprintf("%s %s", m.Cmd, m.Key)
		return []byte(cmd)
	default:
		panic("unknown Command")
	}
}

func parseMessage(rawCmd []byte) (*Message, error) {
	var (
		rawStr = string(rawCmd)
		parts = strings.Split(rawStr, " ")
	)
	
	if len(parts) < 2{
		// respond
		return nil, errors.New("Invalid Command format")
	}
	msg := &Message {
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}
	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, errors.New("Invalid SET Command format")
		}
		msg.Value = []byte(parts[2])
		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println("Invalid integer given with SET Command")
		}
		msg.TTL = time.Duration(ttl)
	}

	return msg, nil
}
