package i2ptcp

import (
	"github.com/eyedeekay/sam3"
	"net"

	manet "github.com/multiformats/go-multiaddr-net"
	ma "github.com/rtradeltd/go-multiaddr"
)

// FromMultiaddrToNetAddr wraps around FromMultiaddrToI2PNetAddr to work with manet.NetCodec
func FromMultiaddrToNetAddr(from ma.Multiaddr) (net.Addr, error) {
	return FromMultiaddrToI2PNetAddr(from)
}

// FromMultiaddrToI2PNetAddr converts a ma.Multiaddr to a sam3.I2PAddr
func FromMultiaddrToI2PNetAddr(from ma.Multiaddr) (sam3.I2PAddr, error) {
	return sam3.NewI2PAddrFromString(from.String())
}

// FromNetAddrToMultiaddr wraps around FromI2PNetAddrToMultiaddr to work with manet.NetCodec
func FromNetAddrToMultiaddr(from net.Addr) (ma.Multiaddr, error) {
	return FromI2PNetAddrToMultiaddr(from.(sam3.I2PAddr))
}

// FromI2PNetAddrToMultiaddr converts a sam3.I2PAddr to a ma.Multiaddr
func FromI2PNetAddrToMultiaddr(from sam3.I2PAddr) (ma.Multiaddr, error) {
	return ma.NewMultiaddr("/garlict/" + from.Base64())
}

func NewGarlicTCPNetCodec() manet.NetCodec {

	var fromNetAddr manet.FromNetAddrFunc
	fromNetAddr = FromNetAddrToMultiaddr

	var toMultiAddr manet.ToNetAddrFunc
	toMultiAddr = FromMultiaddrToNetAddr

	return manet.NetCodec{
		//NetAddrNetworks: ,
		ProtocolName: "garlict",
		// ParseNetAddr parses a net.Addr belonging to this type into a multiaddr
		ParseNetAddr: fromNetAddr,
		// ConvertMultiaddr converts a multiaddr of this type back into a net.Addr
		ConvertMultiaddr: toMultiAddr,
		Protocol:         ma.ProtocolWithName("garlict"),
	}
}
