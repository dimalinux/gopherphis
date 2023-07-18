package jamtis

import (
	"math/big"

	"github.com/dimalinux/gopherphis/util"
)

// edCurveOrder, often called "l", is the prime used by ed25519
var edCurveOrder *big.Int

func init() {
	// python3 -c 'print((2**252 + 27742317777372353535851937790883648493).to_bytes(32, "big").hex())'
	const lHex = "1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed"
	var ok bool
	edCurveOrder, ok = new(big.Int).SetString(lHex, 16)
	if !ok {
		panic("invalid hex constant")
	}
}

// scReduce32 reduces the 32-byte little endian input s by computing and returning
// s mod l, where l is ed25519 curve order prime.
func scReduce32(s32 []byte) []byte {
	scalar := new(big.Int).SetBytes(util.ReverseSlice(s32[:]))
	reduced := util.ReverseSlice(new(big.Int).Mod(scalar, edCurveOrder).Bytes())
	var reduced32 [32]byte
	copy(reduced32[:], reduced) // little endian, so high order byte padding is automatic
	return reduced32[:]
}
