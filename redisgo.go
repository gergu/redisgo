package main

import (
	"os"
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

func handleConnection(conn net.Conn, storage map[string]string) {
	buf := *bufio.NewReader(conn)

	for {
		msg, err := textproto.NewReader(&buf).ReadLine()

		command := strings.SplitN(msg, " ", 2)
		if len(command) < 2 {
			conn.Write([]byte("Invalid instruction. Use following pattern: GET/SET key [value]\n"))
			continue
		}

		if err != nil {
			fmt.Printf("Connection closed.\n")
			break
		}

		switch command[0] {
		case "GET":
			getValue(conn, storage, command[1])
		case "SET":
			setValue(conn, storage, command[1])
		default:
			conn.Write([]byte("Unknown command.\n"))
		}
	}

	conn.Close()
}

func setValue(conn net.Conn, storage map[string]string, command string) {
	instruction := strings.SplitN(command, " ", 2)
	if len(command) < 2 {
		conn.Write([]byte("Invalid SET syntax. Use GET key value\n"))
		return
	}

	storage[instruction[0]] = instruction[1]
	conn.Write([]byte("OK\n"))
}

func getValue(conn net.Conn, storage map[string]string, key string) {
	value := storage[key]
	conn.Write([]byte(value + "\n"))
}

func main() {
	args := os.Args
	PORT := ":6379"

	storage := make(map[string]string)

	if len(args)>1 {
		PORT = ":" + args[1]
	}

	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, storage)
	}
}
