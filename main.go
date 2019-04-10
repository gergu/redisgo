package main

import (
	"bufio"
	"fmt"
	"net"
	"flag"
	"net/textproto"
	"io"
	"strings"
)

type sessionHandler struct {
	conn   io.ReadWriteCloser
	store  map[string]string
}

func (s *sessionHandler) handleConnection() {
	buf := bufio.NewReader(s.conn)

	for {
		msg, err := textproto.NewReader(buf).ReadLine()

		command := strings.SplitN(msg, " ", 2)
		if len(command) < 2 {
			s.conn.Write([]byte("Invalid instruction. Use following pattern: GET/SET key [value]\n"))
			continue
		}

		if err != nil {
			fmt.Printf("Connection closed.\n")
			break
		}

		switch command[0] {
		case "GET":
			s.getValue(command[1])
		case "SET":
			s.setValue(command[1])
		default:
			s.conn.Write([]byte("Unknown command.\n"))
		}
	}

	s.conn.Close()
}

func (s *sessionHandler) setValue(command string) {
	instruction := strings.SplitN(command, " ", 2)
	if len(command) < 2 {
		s.conn.Write([]byte("Invalid SET syntax. Use GET key value\n"))
		return
	}

	s.store[instruction[0]] = instruction[1]
	s.conn.Write([]byte("OK\n"))
}

func (s *sessionHandler) getValue(key string) {
	value := s.store[key]
	s.conn.Write([]byte(value + "\n"))
}

func main() {
	var port string
	var password string
	flag.StringVar(&port, "port", "6379", "port")
	flag.StringVar(&password, "password", "", "password")

	flag.Parse()

	storage := make(map[string]string)

	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		handler := sessionHandler{conn, storage}

		go handler.handleConnection()
	}
}
