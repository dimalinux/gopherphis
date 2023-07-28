package cryptonote

import (
	"bytes"
	"errors"
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// Network is the Monero network type
type Network string

// Monero networks
const (
	Mainnet  Network = "mainnet"
	Stagenet Network = "stagenet"
	Testnet  Network = "testnet"
)

// AddressType is the type of Monero address: Standard or Subaddress
type AddressType string

// Monero address types. We don't support Integrated.
const (
	Standard   AddressType = "standard"
	Subaddress AddressType = "subaddress"
)

// Network prefix byte. The 1st decoded byte of a monero address defines both
// the network (mainnet, stagenet, testnet) and the type of address (standard,
// integrated, and subaddress).
const (
	netPrefixStdAddrMainnet  = 18
	netPrefixSubAddrMainnet  = 42
	netPrefixStdAddrStagenet = 24
	netPrefixSubAddrStagenet = 36
	netPrefixStdAddrTestnet  = 53
	netPrefixSubAddrTestnet  = 63
)

var (
	errAddressNotInitialized    = errors.New("monero address is not initialized")
	errChecksumMismatch         = errors.New("invalid address checksum")
	errInvalidAddressLength     = errors.New("invalid monero address length")
	errInvalidAddressEncoding   = errors.New("invalid monero address encoding")
	errInvalidPrefixGotMainnet  = errors.New("invalid monero address: expected stagenet, got mainnet")
	errInvalidPrefixGotStagenet = errors.New("invalid monero address: expected mainnet, got stagenet")
	errInvalidPrefixGotTestnet  = errors.New("invalid monero address: monero testnet not yet supported")
)

// Address represents a Monero address
type Address struct {
	// decoded is the bytes (prefix, pub spend key, pub view key, checksum) that
	// get base58 encoded. Package private, as it is a semi-arbitrary
	// implementation detail.
	decoded [addressBytesLen]byte
}

// NewAddress converts a string to a monero Address with validation.
func NewAddress(addrStr string, net Network) (*Address, error) {
	addr := new(Address)
	if err := addr.UnmarshalText([]byte(addrStr)); err != nil {
		return nil, err
	}

	return addr, addr.ValidateNet(net)
}

func (a *Address) String() string {
	return addrBytesToBase58(a.decoded[:])
}

// Network returns the Monero network of the address
func (a *Address) Network() Network {
	switch a.decoded[0] {
	case netPrefixStdAddrMainnet, netPrefixSubAddrMainnet:
		return Mainnet
	case netPrefixStdAddrStagenet, netPrefixSubAddrStagenet:
		return Stagenet
	case netPrefixStdAddrTestnet, netPrefixSubAddrTestnet:
		return Testnet
	default:
		// Our methods to deserialize and create Address values all verify
		// that the address byte is valid
		panic("address has invalid network prefix")
	}
}

// Type returns the Address type
func (a *Address) Type() AddressType {
	switch a.decoded[0] {
	case netPrefixStdAddrMainnet, netPrefixStdAddrStagenet, netPrefixStdAddrTestnet:
		return Standard
	case netPrefixSubAddrTestnet, netPrefixSubAddrStagenet, netPrefixSubAddrMainnet:
		return Subaddress
	default:
		// Our methods to deserialize and create Address values all verify
		// that the address byte is valid
		panic("address has invalid network prefix")
	}
}

// validateDecoded ensures that the checksum and network prefix of the address
// are valid. The Network() and Type() methods are not safe to use until
// this base level validation is performed.
func (a *Address) validateDecoded() error {
	checksum := getChecksum(a.decoded[:65])
	if !bytes.Equal(checksum[:], a.decoded[65:69]) {
		return errChecksumMismatch
	}

	netPrefix := a.decoded[0]
	switch netPrefix {
	case netPrefixStdAddrMainnet, netPrefixSubAddrMainnet,
		netPrefixStdAddrStagenet, netPrefixSubAddrStagenet,
		netPrefixStdAddrTestnet, netPrefixSubAddrTestnet:
		// we are good, do nothing
	default:
		return fmt.Errorf("monero address has unknown network prefix %d", netPrefix)
	}

	return nil
}

// Equal returns true if the addresses are identical, otherwise false.
func (a *Address) Equal(b *Address) bool {
	if b == nil {
		return false
	}
	return a.decoded == b.decoded
}

// ValidateNet validates that the monero network matches the passed network.
// This validation can't be performed when decoding JSON, as the environment is
// not known at that time.
func (a *Address) ValidateNet(net Network) error {
	if a == nil || a.decoded == new(Address).decoded {
		return errAddressNotInitialized
	}

	switch a.Network() {
	case Mainnet:
		if net != Mainnet {
			return errInvalidPrefixGotMainnet
		}
	case Stagenet:
		if net != Stagenet {
			return errInvalidPrefixGotStagenet
		}
	case Testnet:
		return errInvalidPrefixGotTestnet
	default:
		panic("unhandled network")
	}

	return nil
}

func getChecksum(data ...[]byte) (result [4]byte) {
	keccak256 := ethcrypto.Keccak256(data...)
	copy(result[:], keccak256[:4])
	return
}

// Address returns the address as bytes for a PublicKeyPair with the given environment (ie. mainnet or stagenet)
func (kp *PublicKeyPair) Address(net Network) *Address {
	address := new(Address)

	var prefix byte
	switch {
	case net == Mainnet && !kp.isSubAddress:
		prefix = netPrefixStdAddrMainnet
	case net == Mainnet && kp.isSubAddress:
		prefix = netPrefixSubAddrMainnet
	case net == Stagenet && !kp.isSubAddress:
		prefix = netPrefixStdAddrStagenet
	case net == Stagenet && kp.isSubAddress:
		prefix = netPrefixSubAddrStagenet
	case net == Testnet && !kp.isSubAddress:
		prefix = netPrefixStdAddrTestnet
	case net == Testnet && kp.isSubAddress:
		prefix = netPrefixSubAddrTestnet
	default:
		panic(fmt.Sprintf("unhandled net %s", net))
	}

	// address encoding is:
	// (network_prefix) + (32-byte public spend key) + (32-byte-byte public view key)
	// + first_4_Bytes(Hash(network_prefix + (32-byte public spend key) + (32-byte public view key)))
	address.decoded[0] = prefix                 // 1-byte network prefix
	copy(address.decoded[1:33], kp.sk.Bytes())  // 32-byte public spend key
	copy(address.decoded[33:65], kp.vk.Bytes()) // 32-byte public view key
	checksum := getChecksum(address.decoded[0:65])
	copy(address.decoded[65:69], checksum[:])

	return address
}
