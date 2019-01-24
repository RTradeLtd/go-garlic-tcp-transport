package i2ptcpconn

import (
	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	conn, err := NewGarlicTCPConnFromOptions(
		SAMHost("127.0.0.1"),
		SAMPort("7656"),
		SAMPass(""),
		KeysPath(""),
		OnlyGarlic(false),
	)
	if err != nil {
		t.Error(err.Error())
	}
	listener, err := conn.ListenI2P()
	if err != nil {
		t.Error(err)
	}
	log.Println(listener.ID())
	log.Println(listener.Base64())
}

func TestGarlicTransportMaStrings(t *testing.T) {
	conn, err := NewGarlicTCPConnFromOptions(
		SAMHost("/ip4/127.0.0.1/"),
		SAMPort("/tcp/7656/"),
		SAMPass(""),
		KeysPath(""),
		OnlyGarlic(false),
	)
	if err != nil {
		t.Error(err.Error())
	}
	listener, err := conn.ListenI2P()
	if err != nil {
		t.Error(err)
	}
	log.Println(listener.ID())
	log.Println(listener.Base64())
}
