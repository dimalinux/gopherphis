package jamtis

import (
	ed25519 "filippo.io/edwards25519"
	"golang.org/x/crypto/blake2b"
)

func blake2bHash(optionalKey []byte, hashSize int, inputs ...any) ([]byte, error) {
	h, err := blake2b.New(hashSize, optionalKey)
	if err != nil {
		return nil, err
	}

	for _, input := range inputs {
		var inputBytes []byte
		switch arg := input.(type) {
		case string:
			inputBytes = []byte(arg)
		case []byte:
			inputBytes = arg
		case *ed25519.Scalar:
			inputBytes = arg.Bytes()
		case *ed25519.Point:
			inputBytes = arg.Bytes()
		default:
			panic("invalid hash input type")
		}

		_, err = h.Write(inputBytes)
		if err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
}
