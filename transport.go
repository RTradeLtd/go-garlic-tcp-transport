package i2ptcp

import (
	"context"
	//"crypto/rand"

	//crypto "github.com/libp2p/go-libp2p-crypto"
	//net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	//peerstore "github.com/libp2p/go-libp2p-peerstore"

	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/rtradeltd/go-garlic-tcp-transport/common"
	"github.com/rtradeltd/go-garlic-tcp-transport/conn"
)

// GarlicTCPTransport is a libp2p interface to an i2p TCP-like tunnel created
// via the SAM bridge
type GarlicTCPTransport struct {
	hostSAM string
	portSAM string
	passSAM string

	keysPath string

	onlyGarlic    bool
	garlicOptions []string
}

func (t GarlicTCPTransport) PrintOptions() []string {
    return t.garlicOptions
}

// CanDial implements transport.CanDial
func (t GarlicTCPTransport) CanDial(m ma.Multiaddr) bool {
	return t.Matches(m)
}

// CanDialI2P is a special CanDial function that only returns true if it's an
// i2p address.
func (t GarlicTCPTransport) CanDialI2P(m ma.Multiaddr) bool {
	return t.MatchesI2P(m)
}

// Matches returns true if the address is a valid garlic TCP multiaddr
func (t *GarlicTCPTransport) Matches(a ma.Multiaddr) bool {
	return i2phelpers.IsValidGarlicMultiAddr(a)
}

// Matches returns true if the address is a valid garlic TCP multiaddr
func (t *GarlicTCPTransport) MatchesI2P(a ma.Multiaddr) bool {
	return i2phelpers.IsValidGarlicMultiAddr(a)
}

// Dial returns a new GarlicConn
func (t GarlicTCPTransport) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	conn, err := i2ptcpconn.NewGarlicTCPConn(t, t.hostSAM, t.portSAM, t.passSAM, t.keysPath, t.onlyGarlic)
	if err != nil {
		return nil, err
	}
	return conn.DialI2P(c, m, p)
}

func (t GarlicTCPTransport) Listen(addr ma.Multiaddr) (tpt.Listener, error) {
	return t.ListenI2P(addr)
}

func (t GarlicTCPTransport) ListenI2P(addr ma.Multiaddr) (*i2ptcpconn.GarlicTCPConn, error) {
	conn, err := i2ptcpconn.NewGarlicTCPConn(t, t.hostSAM, t.portSAM, t.passSAM, t.keysPath, t.onlyGarlic)
	if err != nil {
		return nil, err
	}
	return conn.ListenI2P()
}

// Protocols need only return this I think
func (t GarlicTCPTransport) Protocols() []int {
	return []int{ma.P_GARLIC64}
}

// Proxy always returns false, we're using the SAM bridge to make our requests
func (t GarlicTCPTransport) Proxy() bool {
	return false
}

// NewGarlicTransport initializes a GarlicTransport for libp2p
func NewGarlicTCPTransport(host, port, pass string, keysPath string, onlyGarlic bool, options []string) (tpt.Transport, error) {
	return NewGarlicTCPTransportFromOptions(
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

func NewGarlicTCPTransportFromOptions(opts ...func(*GarlicTCPTransport) error) (*GarlicTCPTransport, error) {
	var g GarlicTCPTransport
	g.hostSAM = "127.0.0.1"
	g.portSAM = "7656"
	g.passSAM = ""
	g.keysPath = ""
	g.onlyGarlic = false
	g.garlicOptions = []string{}
	for _, o := range opts {
		if err := o(&g); err != nil {
			return nil, err
		}
	}
	return &g, nil
}
