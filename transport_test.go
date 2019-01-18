package i2ptcp

import (
	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	transport, err := NewGarlicTCPTransportFromOptions()
	if err != nil {
		t.Error(err.Error())
	}
	listener, err := transport.ListenI2P()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(listener.ID())
	log.Println(listener.Base64())
}
