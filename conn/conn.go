package i2ptcpconn

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"

	"github.com/libp2p/go-stream-muxer"
	"github.com/rtradeltd/go-garlic-tcp-transport/codec"
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
	i2pkeys         *sam3.I2PKeys

	lPrivKey crypto.PrivKey
	lPubKey  crypto.PubKey

	rPubKey crypto.PubKey

	hostSAM string
	portSAM string
	//passSAM string

	keysPath string

	onlyGarlic    bool
	garlicOptions []string
}

func (t *GarlicTCPConn) SAMHost() string {
	st := strings.TrimPrefix(t.hostSAM, "/ip4/")
	stt := strings.TrimPrefix(st, "/ip6/")
	rt := strings.TrimSuffix(stt, "/")
	return rt
}

func (t *GarlicTCPConn) SAMPort() string {
	st := strings.TrimPrefix(t.portSAM, "/tcp/")
	rt := strings.TrimSuffix(st, "/")
	return rt
}

func (t GarlicTCPConn) SAMAddress() string {
	rt := t.SAMHost() + ":" + t.SAMPort()
	fmt.Println(rt)
	return rt
}

func (t GarlicTCPConn) PrintOptions() []string {
	return t.garlicOptions
}

func (g GarlicTCPConn) MaBase64() ma.Multiaddr {
	r, err := i2ptcpcodec.FromI2PNetAddrToMultiaddr(g.i2pkeys.Addr())
	if err != nil {
		panic("Critical address error! There is no way this should have occurred")
	}
	return r
}

func (g GarlicTCPConn) Base32() string {
	return g.i2pkeys.Addr().Base32()
}

func (g GarlicTCPConn) Base64() string {
	return g.i2pkeys.Addr().Base64()
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
	return c.DialI2P(nil, c.RemoteMultiaddr(), c.RemotePeer())
}

// LocalMultiaddr returns the local multiaddr for this connection
func (g GarlicTCPConn) LocalMultiaddr() ma.Multiaddr {
	return g.laddr
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (c GarlicTCPConn) RemoteMultiaddr() ma.Multiaddr {
	return c.MaBase64()
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

//transport keys
func (g GarlicTCPConn) forward(conn *GarlicTCPConn) {
	//var request *http.Request
	var err error
	var client net.Conn
	if client, err = net.Dial("tcp", g.Addr().String()); err != nil {
		panic("Dial failed: %v" + err.Error())
	}
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

// Listen implements a connection, but addr is IGNORED here, it's drawn from the
//transport keys
func (g GarlicTCPConn) Forward(addr ma.Multiaddr) {
	g.ForwardI2P(addr)
}

// ListenI2P is like Listen, but it returns the GarlicTCPConn and doesn't
//require a multiaddr
func (g GarlicTCPConn) ForwardI2P(addr ma.Multiaddr) {
	var err error
	g.laddr = addr
	g.StreamListener, err = g.StreamSession.Listen()
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := g.AcceptI2P()
		if err != nil {
			panic("ERROR: failed to accept listener: %v" + err.Error())
		}
		go g.forward(conn)
	}
}

// Return the net.Addr version of the local Multiaddr
func (g GarlicTCPConn) Addr() net.Addr {
	ra, _ := manet.ToNetAddr(g.Multiaddr())
	return ra
}

// return the local Multiaddr
func (g GarlicTCPConn) Multiaddr() ma.Multiaddr {
	return g.laddr
}

func NewGarlicTCPConn(transport tpt.Transport, host, port, pass string, keysPath string, onlyGarlic bool, options []string) (*GarlicTCPConn, error) {
	return NewGarlicTCPConnFromOptions(
		Transport(transport),
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

// NewGarlicTCPConnFromOptions creates a GarlicTCPConn using function arguments
func NewGarlicTCPConnFromOptions(opts ...func(*GarlicTCPConn) error) (*GarlicTCPConn, error) {
	var g GarlicTCPConn
	g.hostSAM = "127.0.0.1"
	g.portSAM = "7656"
	//g.passSAM = ""
	g.keysPath = ""
	g.onlyGarlic = false
	g.garlicOptions = []string{}
	for _, o := range opts {
		if err := o(&g); err != nil {
			return nil, err
		}
	}
	var err error
	g.SAM, err = sam3.NewSAM(g.SAMAddress())
	if err != nil {
		return nil, err
	}
	g.i2pkeys, err = g.GetI2PKeys()
	if err != nil {
		return nil, err
	}
	g.StreamSession, err = g.SAM.NewStreamSession(i2phelpers.RandTunName(), *g.i2pkeys, g.PrintOptions())
	if err != nil {
		return nil, err
	}
	return &g, nil
}
