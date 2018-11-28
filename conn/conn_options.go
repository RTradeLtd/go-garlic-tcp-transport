package i2ptcpconn

import (
	"fmt"
	tpt "github.com/libp2p/go-libp2p-transport"
	"strconv"
)

type Option func(*GarlicTCPConn) error

//SAMHost sets the host of the SAM Bridge to use
func Transport(t tpt.Transport) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.parentTransport = t
		return nil
	}
}

//SAMHost sets the host of the SAM Bridge to use
func SAMHost(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.hostSAM = s
		return nil
	}
}

//SAMPort sets the port of the SAM bridge to use
func SAMPort(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if val > 0 && val < 65536 {
			c.portSAM = s
			return nil
		}
		return fmt.Errorf("port is %s invalid")
	}
}

//SAMPass sets the password to use when authenticating to the SAM bridge. It's
//ignored for now, and will return an error if it recieves a non-empty string.
func SAMPass(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		if s != "" {
			return fmt.Errorf("SAMPass is unused for now, pass no argument(or empty string). Failing closed.")
		}
		return nil
	}
}

//KeysPath sets the path to the keys, if no keys are present, they will be generated.
func KeysPath(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.keysPath = s
		return nil
	}
}

func OnlyGarlic(b bool) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.onlyGarlic = b
		return nil
	}
}
