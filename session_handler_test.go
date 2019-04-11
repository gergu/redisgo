package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bytes"
)

type ConnectionMock struct {
	bytes.Buffer
}

func (cm *ConnectionMock) Close() error {
	return nil
}

func TestSetValue(t *testing.T) {
	conn := ConnectionMock{}
	storage := make(map[string]string)

	handler := sessionHandler{&conn, storage, "", true}

  handler.setValue("key")
  assert.Equal(t, handler.store["key"], "")
  assert.Equal(t, "Invalid SET syntax. Use GET key value\n", conn.String())
  conn.Reset()

	handler.setValue("key value")
  assert.Equal(t, handler.store["key"], "value")
  assert.Equal(t, "+OK\n", conn.String())
}

func TestGetValue(t *testing.T) {
	conn := ConnectionMock{}
	storage := make(map[string]string)

	storage["key"] = "value"
	handler := sessionHandler{&conn, storage, "", true}

	handler.getValue("key")
  assert.Equal(t, "value\n", conn.String())
}


func TestAuth(t *testing.T) {
	conn := ConnectionMock{}
	storage := make(map[string]string)

	handler := sessionHandler{&conn, storage, "password", false}

	// Invalid password

	handler.auth("invalid_password")

  assert.Equal(t, "-ERR invalid password\n", conn.String())

  conn.Reset()

  // Valid password

  handler.auth("password")
  assert.Equal(t, "+OK\n", conn.String())
  assert.Equal(t, true, handler.authenticated)
}


func TestHandleCommand(t *testing.T) {
	conn := ConnectionMock{}
	storage := make(map[string]string)

	handler := sessionHandler{&conn, storage, "password", false}
	handler.handleCommand("GET key")
  assert.Equal(t, "-NOAUTH Authentication required.\n", conn.String())
  conn.Reset()

  handler.password = ""
  handler.authenticated = true

  handler.handleCommand("GET")
  assert.Equal(t, "Invalid instruction. Use following pattern: GET/SET key [value]\n", conn.String())
  conn.Reset()

  handler.handleCommand("UNKNOWN parameter param")
  assert.Equal(t, "Unknown command.\n", conn.String())
  conn.Reset()
}
