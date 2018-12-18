package i2phelpers

import (
	"fmt"
	"github.com/eyedeekay/sam3"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	ma "github.com/rtradeltd/go-multiaddr"
)

// IsValidGarlicMultiAddr is used to validate that a multiaddr
// is representing a I2P garlic service
func IsValidGarlicMultiAddr(a ma.Multiaddr) bool {
	if len(a.Protocols()) < 2 {
		return false
	}

	// check for correct network type
	if a.Protocols()[0].Name != "garlic64" {
		fmt.Println("Protocol != garlic64")
		return false
	}

	// split into garlic64 address
	addr, err := a.ValueForProtocol(ma.P_GARLIC64)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//kinda crude, but if it's bigger than this it's at least possible that
	//it's a valid kind of i2p address.
	if len(addr) == 516 {
		fmt.Println(addr)
		return false
	}

	return true
}

// RandTunName generates a random tunnel names to avoid collisions
func RandTunName() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyz"[rand.Intn(len("abcdefghijklmnopqrstuvwxyz"))]
	}
	return string(b)
}

// LoadKeys loads keys into our keys from files in the keys directory
func LoadKeys(keysPath string) (*sam3.I2PKeys, error) {
	absPath, err := filepath.EvalSymlinks(keysPath)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(absPath, ".i2pkeys") {
		file, err := os.Open(absPath)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		keys, err := sam3.LoadKeysIncompat(file)
		if err != nil {
			return nil, err
		}
		return &keys, nil
	}

	return CreateEepServiceKey()
}

func CreateEepServiceKey() (*sam3.I2PKeys, error) {
	sam, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	k, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	return &k, err
}

func EepServiceMultiAddr() (*ma.Multiaddr, error) {
	k, err := CreateEepServiceKey()
	if err != nil {
		return nil, err
	}
	m, err := ma.NewMultiaddr("/garlic64/" + k.Addr().Base64() + ":80")
	if err != nil {
		return nil, err
	}
	return &m, nil
}
