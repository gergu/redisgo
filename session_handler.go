package main

import (
  "bufio"
  "fmt"
  "net/textproto"
  "io"
  "strings"
)

type sessionHandler struct {
  conn   io.ReadWriteCloser
  store  map[string]string
  password string
  authenticated bool
}

func (s *sessionHandler) handleConnection() {
  buf := bufio.NewReader(s.conn)

  defer s.conn.Close()

  for {
    msg, err := textproto.NewReader(buf).ReadLine()

    if err != nil {
      fmt.Printf("Connection closed.\n")
      break
    }

    s.handleCommand(msg)
  }
}

func (s *sessionHandler) handleCommand(msg string) {
  command := strings.SplitN(msg, " ", 2)
  if len(command) < 2 {
    s.conn.Write([]byte("Invalid instruction. Use following pattern: GET/SET key [value]\n"))
    return
  }

  if command[0] == "AUTH" {
    s.auth(command[1])
    return
  }

  if !s.authenticated {
    s.conn.Write([]byte("-NOAUTH Authentication required.\n"))
    return
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

func (s *sessionHandler) setValue(command string) {
  instruction := strings.SplitN(command, " ", 2)
  if len(instruction) < 2 {
    s.conn.Write([]byte("Invalid SET syntax. Use GET key value\n"))
    return
  }

  s.store[instruction[0]] = instruction[1]
  s.conn.Write([]byte("+OK\n"))
}

func (s *sessionHandler) getValue(key string) {
  value := s.store[key]
  s.conn.Write([]byte(value + "\n"))
}

func (s *sessionHandler) initAuth() {
  if s.password != "" {
    s.authenticated = false
  }
}

func (s *sessionHandler) auth(passwd string) {
  if passwd == s.password {
    s.authenticated = true
    s.conn.Write([]byte("+OK\n"))
  } else {
    s.conn.Write([]byte("-ERR invalid password\n"))
  }
}
