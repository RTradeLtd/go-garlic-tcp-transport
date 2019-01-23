package i2ptcp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Option is a functional argument
type Option func(*GarlicTCPTransport) error

//SAMHost sets the host of the SAM Bridge to use
func SAMHost(s string) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		st := ""
		if !strings.Contains("/ip4/", s) || !strings.Contains("/ip6/", s) {
			ip := net.ParseIP(s)
			if ip.To4() != nil {
				st = "/ip4/" + s + "/"
			}
			if ip.To16() != nil {
				st = "/ip6/" + s + "/"
			}
		}
		c.hostSAM = st
		return nil
	}
}

//SAMPort sets the port of the SAM bridge to use
func SAMPort(s string) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		st := strings.Replace("/", "", strings.Replace("/tcp/", "", s, -1), -1)
		val, err := strconv.Atoi(st)
		if err != nil {
			return err
		}
		if val > 0 && val < 65536 {
			c.portSAM = "/tcp/" + st + "/"
			return nil
		}
		return fmt.Errorf("port is %s invalid", s)
	}
}

//SAMPass sets the password to use when authenticating to the SAM bridge. It's
//ignored for now, and will return an error if it recieves a non-empty string.
func SAMPass(s string) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		if s != "" {
			return fmt.Errorf("SAMPass is unused for now, pass no argument(or empty string). Failing closed.")
		}
		return nil
	}
}

//KeysPath sets the path to the keys, if no keys are present, they will be generated.
func KeysPath(s string) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		c.keysPath = s
		return nil
	}
}

func OnlyGarlic(b bool) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		c.onlyGarlic = b
		return nil
	}
}

func GarlicOptions(s []string) func(*GarlicTCPTransport) error {
	return func(c *GarlicTCPTransport) error {
		for _, v := range s {
			c.garlicOptions = append(c.garlicOptions, v)
		}
		return nil
	}
}
