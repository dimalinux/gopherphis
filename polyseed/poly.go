package polyseed

const (
	numChecksumWords = 1

	keyBitsPerWord = 10
	dataWords      = NumSeedWords - numChecksumWords

	bitsPerByte = 8
)

type poly struct {
	coeff [NumSeedWords]uint16
}

var mul2Table = [8]uint16{
	5, 7, 1, 3, 13, 15, 9, 11,
}

func uint16Mul2(x uint16) uint16 {
	if x < 1024 {
		return 2 * x
	}
	return mul2Table[x%8] + 16*((x-1024)/8)
}

//func newGfPoly

func (p *poly) calcChecksum() uint16 {
	// Horner's method at x = 2
	sum := p.coeff[NumSeedWords-1]
	for i := NumSeedWords - 2; i >= 0; i-- {
		sum = uint16Mul2(sum) ^ p.coeff[i]
	}
	return sum
}

func (p *poly) ValidChecksum() bool {
	return p.calcChecksum() == 0
}

func dataToPoly(data *SeedData) *poly {

	extraVal := (data.features << DateBits) | data.birthday
	extraBits := uint(featureBits + DateBits)

	var wordBits, wordVal, secretIdx uint
	secretVal := uint(data.secret[secretIdx])
	secretBits := uint(bitsPerByte)
	seedRemBits := uint(NumSecretBits - bitsPerByte)

	poly := &poly{}

	for i := 0; i < dataWords; i++ {
		for wordBits < keyBitsPerWord {
			if secretBits == 0 {
				secretIdx++
				secretBits = min(seedRemBits, bitsPerByte)
				secretVal = uint(data.secret[secretIdx])
				seedRemBits -= secretBits
			}
			chunkBits := min(secretBits, keyBitsPerWord-wordBits)
			secretBits -= chunkBits
			wordBits += chunkBits
			wordVal <<= chunkBits
			wordVal |= (secretVal >> secretBits) & ((1 << chunkBits) - 1)
		}
		wordVal <<= 1
		extraBits--
		wordVal |= (extraVal >> extraBits) & 1
		poly.coeff[numChecksumWords+i] = uint16(wordVal)
		wordVal = 0
		wordBits = 0
	}

	poly.coeff[0] = poly.calcChecksum()
	poly.coeff[numChecksumWords] ^= uint16(defaultCoin)

	if seedRemBits != 0 {
		panic("seed_rem_bits is not zero")
	}

	if secretBits != 0 {
		panic("secret_bits are not zero")
	}

	if extraBits != 0 {
		panic("extra bits are not zero")
	}

	return poly
}

func (p *poly) ToSeedData() *SeedData {
	data := &SeedData{}

	var extraVal, extraBits, wordBits, secretIdx, secretBits, seedBits uint

	// First word only has checksum bits, so we skip it
	for i := 1; i < NumSeedWords; i++ {
		wordIndex := uint(p.coeff[i]) // the 11-bit index from a language word list

		extraVal <<= 1
		extraVal |= wordIndex & 1
		wordIndex >>= 1
		extraBits++

		for wordBits = 10; wordBits > 0; {
			if secretBits == bitsPerByte {
				secretIdx++
				seedBits += secretBits
				secretBits = 0
			}
			chunkBits := min(wordBits, bitsPerByte-secretBits)
			wordBits -= chunkBits
			chunkMask := (uint(1) << chunkBits) - 1
			if chunkBits < bitsPerByte {
				data.secret[secretIdx] <<= chunkBits
			}
			data.secret[secretIdx] |= uint8((wordIndex >> wordBits) & chunkMask)
			secretBits += chunkBits
		}
	}

	seedBits += secretBits

	if wordBits != 0 {
		panic("word_bits is not zero")
	}

	if seedBits != NumSecretBits {
		panic("seed_bits is not SecretBits")
	}

	if extraBits != featureBits+DateBits {
		panic("extra_bits has wrong value")
	}

	data.birthday = extraVal & DateMask
	data.features = extraVal >> DateBits

	return data
}
