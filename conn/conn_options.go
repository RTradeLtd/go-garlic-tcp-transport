package i2ptcpconn

import (
	"fmt"

	//peer "github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-transport"
)

// Option is a functional argument to the connection constructor
type Option func(*GarlicTCPConn) error

//Transport sets the parent transport of the connection.
func Transport(t tpt.Transport) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.parentTransport = t
		return nil
	}
}

//SAMPass sets the password to use when authenticating to the SAM bridge. It's
//ignored for now, and will return an error if it recieves a non-empty string.
func SAMPass(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		if s != "" {
			return fmt.Errorf("SAMPass is unused for now, pass no argument(or empty string). Failing closed")
		}
		return nil
	}
}

//OnlyGarlic indicates that this connection will only be used to serve anonymous
//connections. It does nothing but indicate that for now.
func OnlyGarlic(b bool) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.onlyGarlic = b
		return nil
	}
}

// GarlicOptions is a slice of string-formatted options to pass to the SAM API.
func GarlicOptions(s []string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		for _, v := range s {
			c.garlicOptions = append(c.garlicOptions, v)
		}
		return nil
	}
}

/*
func LocalPeerID(p peer.ID) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.id = p
		return nil
	}
}
*/
