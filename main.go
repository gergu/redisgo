package main

import (
	"net"
	"flag"
)

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

		handler := sessionHandler{conn, storage, password, true}
		handler.initAuth()

		go handler.handleConnection()
	}
}
