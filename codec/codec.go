package i2ptcpcodec

import (
	//	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/geti2p64"
	"github.com/eyedeekay/sam3/i2pkeys"
	"net"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
)

// FromMultiaddrToNetAddr wraps around FromMultiaddrToI2PNetAddr to work with manet.NetCodec, requires a full base64 to work
func FromMultiaddrToNetAddr(from ma.Multiaddr) (net.Addr, error) {
	return FromMultiaddrToI2PNetAddr(from)
}

// FromMultiaddrToI2PNetAddr converts a ma.Multiaddr to a sam3.I2PAddr, requires a full base64 to work
func FromMultiaddrToI2PNetAddr(from ma.Multiaddr) (i2pkeys.I2PAddr, error) {
	if strings.HasSuffix(from.String(), ".i2p") {
		final, err := lookup.Lookup(from.String())
		if err == nil {
			return i2pkeys.NewI2PAddrFromString(final)
		}
	}
	return i2pkeys.NewI2PAddrFromString(from.String())
}

// FromNetAddrToMultiaddr wraps around FromI2PNetAddrToMultiaddr to work with manet.NetCodec
func FromNetAddrToMultiaddr(from net.Addr) (ma.Multiaddr, error) {
	return FromI2PNetAddrToMultiaddr(from.(i2pkeys.I2PAddr))
}

// FromI2PNetAddrToMultiaddr converts a sam3.I2PAddr to a ma.Multiaddr
func FromI2PNetAddrToMultiaddr(from i2pkeys.I2PAddr) (ma.Multiaddr, error) {
	return ma.NewMultiaddr("/garlic64/" + from.Base64() + "/garlic32/" + from.Base32())
}
