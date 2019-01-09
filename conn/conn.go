package i2ptcpconn

import (
	"context"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	"net"

	"github.com/libp2p/go-stream-muxer"
	"github.com/rtradeltd/go-garlic-tcp-transport/common"
	"github.com/rtradeltd/sam3"
)

// GarlicTCPConn implements a Conn interface
type GarlicTCPConn struct {
	*sam3.SAMConn
	*sam3.SAM
	*sam3.StreamSession
	*sam3.StreamListener
	parentTransport tpt.Transport
	laddr           ma.Multiaddr
	raddr           ma.Multiaddr

	lPrivKey crypto.PrivKey
	lPubKey  crypto.PubKey

	rPubKey crypto.PubKey

	hostSAM string
	portSAM string
	//passSAM string

	keysPath string

	onlyGarlic bool
}

// Tranpsort returns the GarlicTCPTransport to which the GarlicTCPConn belongs
func (g GarlicTCPConn) Transport() tpt.Transport {
	return g.parentTransport
}

// IsClosed says a connection is closed if g.StreamSession is nil because
// Close() nils it if it works. Might need to re-visit that.
func (g GarlicTCPConn) IsClosed() bool {
	if g.StreamSession == nil {
		return true
	}
	return false
}

// AcceptStream lets us streammux
func (c GarlicTCPConn) AcceptStream() (streammux.Stream, error) {
	return c.AcceptI2P()
}

//
func (g GarlicTCPConn) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	return g.DialI2P(c, m, p)
}

func (g GarlicTCPConn) DialI2P(c context.Context, m ma.Multiaddr, p peer.ID) (*GarlicTCPConn, error) {
	var err error
	g.SAMConn, err = g.StreamSession.DialContextI2P(c, "", m.String())
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// OpenStream lets us streammux
func (c GarlicTCPConn) OpenStream() (streammux.Stream, error) {
	return c.DialI2P(nil, c.raddr, c.RemotePeer())
}

// LocalMultiaddr returns the local multiaddr for this connection
func (g GarlicTCPConn) LocalMultiaddr() ma.Multiaddr {
	return g.laddr
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (c GarlicTCPConn) RemoteMultiaddr() ma.Multiaddr {
	return c.raddr
}

// LocalPrivateKey returns the local private key used for the peer.ID
func (c GarlicTCPConn) LocalPrivateKey() crypto.PrivKey {
	return c.lPrivKey
}

// RemotePeer returns the remote peer.ID used for IPFS
func (c GarlicTCPConn) RemotePeer() peer.ID {
	rpeer, err := peer.IDFromPublicKey(c.RemotePublicKey())
	if err != nil {
		panic(err)
	}
	return rpeer
}

//RemotePublicKey returns the remote public key used for the peer.ID
func (c GarlicTCPConn) RemotePublicKey() crypto.PubKey {
	return c.rPubKey
}

// LocalPeer returns the local peer.ID used for IPFS
func (c GarlicTCPConn) LocalPeer() peer.ID {
	lpeer, err := peer.IDFromPrivateKey(c.LocalPrivateKey())
	if err != nil {
		panic(err)
	}
	return lpeer
}

// Close ends a SAM session associated with a transport
func (g GarlicTCPConn) Close() error {
	err := g.StreamSession.Close()
	if err == nil {
		g.StreamSession = nil
	}
	return err
}

// Reset lets us streammux, I need to re-examine how to implement it.
func (g GarlicTCPConn) Reset() error {
	return g.Close()
}

// GetI2PKeys loads the i2p address keys and returns them.
func (g GarlicTCPConn) GetI2PKeys() (*sam3.I2PKeys, error) {
	return i2phelpers.LoadKeys(g.keysPath)
}

// Accept implements a listener
func (g GarlicTCPConn) Accept() (tpt.Conn, error) {
	return g.AcceptI2P()
}

// AcceptI2P helps with Accept
func (g GarlicTCPConn) AcceptI2P() (*GarlicTCPConn, error) {
	var err error
	g.SAMConn, err = g.StreamListener.AcceptI2P()
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// Listen implements a listener
func (g GarlicTCPConn) Listen() (tpt.Conn, error) {
	return g.ListenI2P()
}

// ListenI2P helps with Listen
func (g GarlicTCPConn) ListenI2P() (*GarlicTCPConn, error) {
	var err error
	g.StreamListener, err = g.StreamSession.Listen()
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// Return the net.Addr version of the local Multiaddr
func (g GarlicTCPConn) Addr() net.Addr {
	return g.StreamListener.Addr()
}

// return the local Multiaddr
func (g GarlicTCPConn) Multiaddr() ma.Multiaddr {
	return g.laddr
}

func NewGarlicTCPConn(transport tpt.Transport, host, port, pass string, keysPath string, onlyGarlic bool) (*GarlicTCPConn, error) {
	return NewGarlicTCPConnFromOptions(
		Transport(transport),
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
	)
}

// NewGarlicTCPConnFromOptions creates a GarlicTCPConn using function arguments
func NewGarlicTCPConnFromOptions(opts ...func(*GarlicTCPConn) error) (*GarlicTCPConn, error) {
	var g GarlicTCPConn
	for _, o := range opts {
		if err := o(&g); err != nil {
			return nil, err
		}
	}
	var err error
	g.SAM, err = sam3.NewSAM(g.hostSAM + ":" + g.portSAM)
	if err != nil {
		return nil, err
	}
	i2pkeys, err := g.GetI2PKeys()
	if err != nil {
		return nil, err
	}
	g.StreamSession, err = g.SAM.NewStreamSession(i2phelpers.RandTunName(), *i2pkeys, sam3.Options_Medium)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
