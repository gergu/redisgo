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

  assert.Equal(t, handler.store["key"], "value", "they should be equal")
}
