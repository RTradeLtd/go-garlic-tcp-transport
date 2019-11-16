package i2ptcpconn

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/libp2p/go-libp2p-core/mux"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/RTradeLtd/go-garlic-tcp-transport/codec"
	"github.com/RTradeLtd/go-garlic-tcp-transport/common"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

// GarlicTCPConn implements a Conn interface
type GarlicTCPConn struct {
	*sam3.SAMConn
	*sam3.SAM
	*sam3.StreamSession
	*sam3.StreamListener
	i2pkeys.I2PKeys
	network.ConnSecurity

	parentTransport tpt.Transport

	onlyGarlic    bool
	garlicOptions []string
}

var gc tpt.CapableConn = &GarlicTCPConn{}

func (t *GarlicTCPConn) keysPath() string {
	if t.parentTransport != nil {
		x := reflect.TypeOf(t.parentTransport)
		if x.String() == "i2ptcp.GarlicTCPTransport" {
			if _, b := x.FieldByName("keysPath"); b {
				return reflect.ValueOf(t.parentTransport).FieldByName("keysPath").String()
			}
		}
	}
	return "127.0.0.1"
}

// SAMHost returns the IP address of the configured SAM bridge
func (t *GarlicTCPConn) SAMHost() string {
	if t.parentTransport != nil {
		x := reflect.TypeOf(t.parentTransport)
		if x.String() == "i2ptcp.GarlicTCPTransport" {
			if _, b := x.FieldByName("HostSAM"); b {
				return reflect.ValueOf(t.parentTransport).FieldByName("HostSAM").String()
			}
		}
	}
	return "127.0.0.1"
}

// SAMPort returns the Port of the configured SAM bridge
func (t *GarlicTCPConn) SAMPort() string {
	if t.parentTransport != nil {
		x := reflect.TypeOf(t.parentTransport)
		if x.String() == "i2ptcp.GarlicTCPTransport" {
			if _, b := x.FieldByName("PortSAM"); b {
				return reflect.ValueOf(t.parentTransport).FieldByName("PortSAM").String()
			}
		}
	}
	return "7656"
}

// SAMAddress combines them and returns a full address.
func (t *GarlicTCPConn) SAMAddress() string {
	rt := t.SAMHost() + ":" + t.SAMPort()
	fmt.Println(rt)
	return rt
}

func (t *GarlicTCPConn) i2pkey() i2pkeys.I2PKeys {
	if t.I2PKeys.String() == "" {
		t.I2PKeys, _ = t.GetI2PKeys()
	}
	return t.I2PKeys
}

// PrintOptions returns the options passed to the SAM bridge as a slice of
// strings.
func (t *GarlicTCPConn) PrintOptions() []string {
	return t.garlicOptions
}

// MaBase64 gives us a multiaddr by converting an I2PAddr
func (t *GarlicTCPConn) MA() ma.Multiaddr {
	r, err := i2ptcpcodec.FromI2PNetAddrToMultiaddr(t.i2pkey().Addr())
	if err != nil {
		panic("Critical address error! There is no way this should have occurred")
	}
	return r
}

// RemoteMA gives us a multiaddr for the remote peer
func (t *GarlicTCPConn) RemoteMA() ma.Multiaddr {
	r, err := i2ptcpcodec.FromI2PNetAddrToMultiaddr(t.SAMConn.RemoteAddr().(i2pkeys.I2PAddr))
	if err != nil {
		panic("Critical address error! There is no way this should have occurred")
	}
	return r
}

// Base32 returns the remotely-accessible base32 address of the gateway over i2p
// this is the one you want to use to visit it in the browser.
func (t *GarlicTCPConn) Base32() string {
	return t.i2pkey().Addr().Base32()
}

// Base64 returns the remotely-accessible base64 address of the gateway over I2P
func (t *GarlicTCPConn) Base64() string {
	return t.i2pkey().Addr().Base64()
}

// Transport returns the GarlicTCPTransport to which the GarlicTCPConn belongs
func (t *GarlicTCPConn) Transport() tpt.Transport {
	return t.parentTransport
}

// IsClosed says a connection is closed if t.StreamSession is nil because
// Close() nils it if it works. Might need to re-visit that.
func (t *GarlicTCPConn) IsClosed() bool {
	if t.StreamSession == nil {
		return true
	}
	return false
}

// AcceptStream lets us streammux
func (t *GarlicTCPConn) AcceptStream() (mux.MuxedStream, error) {
	return t.AcceptI2P()
}

// Dial dials an I2P client connection to an i2p hidden service using a garlic64
// multiaddr and returns a tpt.CapableConn
func (t *GarlicTCPConn) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.CapableConn, error) {
	return t.DialI2P(c, m, p)
}

// DialI2P helps with Dial and returns a GarlicTCPConn
func (t *GarlicTCPConn) DialI2P(c context.Context, m ma.Multiaddr, p peer.ID) (*GarlicTCPConn, error) {
	var err error
	t.SAMConn, err = t.StreamSession.DialContextI2P(c, "", m.String())
	if err != nil {
		return nil, err
	}
	return t, nil
}

// OpenStream lets us streammux
func (t *GarlicTCPConn) OpenStream() (mux.MuxedStream, error) {
	return t.DialI2P(nil, t.RemoteMultiaddr(), t.RemotePeer())
}

// LocalMultiaddr returns the local multiaddr for this connection
func (t *GarlicTCPConn) LocalMultiaddr() ma.Multiaddr {
	return t.MA()
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (t *GarlicTCPConn) RemoteMultiaddr() ma.Multiaddr {
	return t.RemoteMA()
}

// Close ends a SAM session associated with a transport
func (t *GarlicTCPConn) Close() error {
	err := t.StreamSession.Close()
	if err == nil {
		t.StreamSession = nil
	}
	return err
}

// Reset lets us streammux, I need to re-examine how to implement it.
func (t *GarlicTCPConn) Reset() error {
	return t.Close()
}

// GetI2PKeys loads the i2p address keys and returns them.
func (t *GarlicTCPConn) GetI2PKeys() (i2pkeys.I2PKeys, error) {
	if t.I2PKeys.String() == "" {
		return i2phelpers.LoadKeys(t.keysPath())
	}
	return t.I2PKeys, nil
}

// Accept implements a listener
func (t *GarlicTCPConn) Accept() (tpt.CapableConn, error) {
	return t.AcceptI2P()
}

// AcceptI2P helps with Accept
func (t *GarlicTCPConn) AcceptI2P() (*GarlicTCPConn, error) {
	var err error
	t.SAMConn, err = t.StreamListener.AcceptI2P()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Listen implements a listener
func (t *GarlicTCPConn) Listen() (tpt.CapableConn, error) {
	return t.ListenI2P()
}

// ListenI2P helps with Listen
func (t *GarlicTCPConn) ListenI2P() (*GarlicTCPConn, error) {
	var err error
	t.StreamListener, err = t.StreamSession.Listen()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Addr returns the net.Addr version of the local Multiaddr
func (t *GarlicTCPConn) Addr() net.Addr {
	return t.i2pkey().Addr()
}

// Multiaddr returns the local Multiaddr
func (t *GarlicTCPConn) Multiaddr() ma.Multiaddr {
	return t.LocalMultiaddr()
}

// Stat returns the local Multiaddr
func (t *GarlicTCPConn) Stat() network.Stat {
	var dir network.Direction
	if t.StreamListener == nil {
		dir = network.DirOutbound
	} else if t.SAMConn != nil {
		dir = network.DirInbound
	} else {
		dir = network.DirUnknown
	}
	return network.Stat{
		Direction: dir,
	}
}

// NewGarlicTCPConn creates an I2P Connection struct from a fixed list of arguments
func NewGarlicTCPConn(transport tpt.Transport, onlyGarlic bool, options []string) (*GarlicTCPConn, error) {
	return NewGarlicTCPConnFromOptions(
		Transport(transport),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

// NewGarlicTCPConnFromOptions creates a GarlicTCPConn using function arguments
func NewGarlicTCPConnFromOptions(opts ...func(*GarlicTCPConn) error) (*GarlicTCPConn, error) {
	var t GarlicTCPConn
	t.onlyGarlic = false
	t.garlicOptions = []string{}
	t.parentTransport = nil
	for _, o := range opts {
		if err := o(&t); err != nil {
			return nil, err
		}
	}
	var err error
	if t.parentTransport == nil {
		return nil, fmt.Errorf("Parent transport must be set")
	}
	t.SAM, err = sam3.NewSAM(t.SAMAddress())
	if err != nil {
		return nil, err
	}
	t.I2PKeys, err = t.GetI2PKeys()
	if err != nil {
		return nil, err
	}
	t.StreamSession, err = t.SAM.NewStreamSession(i2phelpers.RandTunName(), t.I2PKeys, t.PrintOptions())
	if err != nil {
		return nil, err
	}
	return &t, nil
}
