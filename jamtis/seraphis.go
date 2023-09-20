package jamtis

import (
	ed25519 "filippo.io/edwards25519"

	"github.com/dimalinux/gopherphis/mcrypto"
)

func getUPoint() *ed25519.Point {
	U, err := new(ed25519.Point).SetBytes([]byte{
		0x10, 0x94, 0x8b, 0x00, 0xd2, 0xde, 0x50, 0xb5,
		0x76, 0x99, 0x8c, 0x11, 0xe8, 0x3c, 0x59, 0xa7,
		0x96, 0x84, 0xd2, 0x5c, 0x9f, 0x8a, 0x0d, 0xc6,
		0x86, 0x45, 0x70, 0xd7, 0x97, 0xb9, 0xc1, 0x6e,
	})
	if err != nil {
		panic(err) // unreachable
	}
	return U
}

func getXPoint() *ed25519.Point {
	X, err := new(ed25519.Point).SetBytes([]byte{
		0xa4, 0xfb, 0x43, 0xca, 0x69, 0x5e, 0x12, 0x99,
		0x88, 0x02, 0xa2, 0x0a, 0x15, 0x8f, 0x12, 0xea,
		0x79, 0x47, 0x4f, 0xb9, 0x01, 0x21, 0x16, 0x95,
		0x6a, 0x69, 0x76, 0x7c, 0x4d, 0x41, 0x11, 0x0f,
	})
	if err != nil {
		panic(err) // unreachable
	}
	return X
}

// GenSeraphisSpendKey returns the ed25519 calculation viewBalance*X + masterKey*U
func GenSeraphisSpendKey(viewBalanceKey []byte, masterKey []byte) ([]byte, error) {
	// Normally, you would call SetBytesWithClamping to get the reduced values,
	// but it has additional bit modifications that are incompatible with
	// Monero, so we do the reduce ourselves and call SetCanonicalBytes instead.
	viewBalanceKey = mcrypto.ScReduce32(viewBalanceKey)
	masterKey = mcrypto.ScReduce32(masterKey)

	mkScalar, err := new(ed25519.Scalar).SetCanonicalBytes(masterKey)
	if err != nil {
		return nil, err
	}

	vbScalar, err := new(ed25519.Scalar).SetCanonicalBytes(viewBalanceKey)
	if err != nil {
		return nil, err
	}

	vbX := new(ed25519.Point).ScalarMult(vbScalar, getXPoint())
	mkU := new(ed25519.Point).ScalarMult(mkScalar, getUPoint())

	return new(ed25519.Point).Add(vbX, mkU).Bytes(), nil
}
