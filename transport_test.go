package i2ptcp

import (
	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	transport, err := NewGarlicTCPTransportFromOptions(
		SAMHost("127.0.0.1"),
		SAMPort("7656"),
		SAMPass(""),
		KeysPath(""),
		OnlyGarlic(false),
	)
	if err != nil {
		t.Error(err.Error())
	}
	listener, err := transport.ListenI2P()
	if err != nil {
		t.Error(err)
	}
	log.Println(listener.ID())
	log.Println(listener.Base64())
}
