package i2ptcp

import (
	"github.com/rtradeltd/go-garlic-tcp-transport/common"
	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	transport, err := NewGarlicTCPTransportFromOptions()
	if err != nil {
		t.Error(err.Error())
	}
	maserver, err := i2phelpers.EepServiceMultiAddr()
	if err != nil {
		t.Error(err.Error())
	}
	listener, err := transport.ListenI2P(*maserver)
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(listener.ID())
	log.Println(maserver)
}
