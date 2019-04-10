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

	handler.setValue("key value")

  assert.Equal(t, handler.store["key"], "value")
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
