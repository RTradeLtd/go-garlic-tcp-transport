package i2phelpers

import (
	"fmt"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	// DefaultPathName is the default config dir name
	KeysPathName = ".ipfs"
	// DefaultPathRoot is the path to the default config dir location.
	KeysPathRoot = "~/" + KeysPathName
	// EnvDir is the environment variable used to change the path root.
	EnvDir = "KEYS_PATH"
)

func Path(filename, extension string) (string, error) {
	dir, err := PathRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename, extension), nil
}

// PathRoot returns the default configuration root directory
func PathRoot() (string, error) {
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(KeysPathRoot)
	}
	return dir, err
}

// IsValidGarlicMultiAddr is used to validate that a multiaddr
// is representing a I2P garlic service
func IsValidGarlicMultiAddr(a ma.Multiaddr) bool {
	if len(a.Protocols()) < 2 {
		return false
	}

	// check for correct network type
	if a.Protocols()[0].Name != "garlic32" {
		fmt.Println("Protocol != garlic32")
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
	if len(addr) < 51 {
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

func isValidExtension(ext string) bool {
	switch ext {
	case
		".i2pkeys",
		".dat":
		return true
	}
	return false
}

// LoadKeys loads keys into our keys from files in the keys directory
func LoadKeys(keysPath string) (i2pkeys.I2PKeys, error) {
	title := filepath.Base(keysPath)
	extension := strings.ToLower(filepath.Ext(title))
	realPath, err := Path(title, extension)
	if err != nil {
		return i2pkeys.I2PKeys{}, err
	}
	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		return CreateEepServiceKey()
	}
	if isValidExtension(extension) {
		file, err := os.Open(realPath)
		defer file.Close()
		if err != nil {
			return i2pkeys.I2PKeys{}, err
		}
		keys, err := i2pkeys.LoadKeysIncompat(file)
		if err != nil {
			return i2pkeys.I2PKeys{}, err
		}
		return keys, nil
	}
	return i2pkeys.I2PKeys{}, fmt.Errorf("Not permitted file extension was encountered.")
}

func CreateEepServiceKey() (i2pkeys.I2PKeys, error) {
	sam, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return i2pkeys.I2PKeys{}, err
	}
	defer sam.Close()
	k, err := sam.NewKeys()
	if err != nil {
		return i2pkeys.I2PKeys{}, err
	}
	return k, err
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
