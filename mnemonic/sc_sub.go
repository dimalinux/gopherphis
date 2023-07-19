package mnemonic

import (
	"encoding/binary"
)

//
//nolint:lll
// This file is a Go port of the following code:
// https://github.com/monero-project/monero/blob/v0.18.2.2/src/crypto/crypto_ops_builder/ref10CommentedCombined/sc_sub.xmr.c
//

func load3(in []byte) int64 {
	return int64(uint64(in[0]) | uint64(in[1])<<8 | uint64(in[2])<<16)
}

func load4(in []byte) int64 {
	return int64(binary.LittleEndian.Uint32(in))
}

func scSub(a []byte, b []byte) []byte {
	var s [32]byte // returned result

	a0 := 2097151 & load3(a)
	a1 := 2097151 & (load4(a[2:]) >> 5)
	a2 := 2097151 & (load3(a[5:]) >> 2)
	a3 := 2097151 & (load4(a[7:]) >> 7)
	a4 := 2097151 & (load4(a[10:]) >> 4)
	a5 := 2097151 & (load3(a[13:]) >> 1)
	a6 := 2097151 & (load4(a[15:]) >> 6)
	a7 := 2097151 & (load3(a[18:]) >> 3)
	a8 := 2097151 & load3(a[21:])
	a9 := 2097151 & (load4(a[23:]) >> 5)
	a10 := 2097151 & (load3(a[26:]) >> 2)
	a11 := load4(a[28:]) >> 7

	b0 := 2097151 & load3(b)
	b1 := 2097151 & (load4(b[2:]) >> 5)
	b2 := 2097151 & (load3(b[5:]) >> 2)
	b3 := 2097151 & (load4(b[7:]) >> 7)
	b4 := 2097151 & (load4(b[10:]) >> 4)
	b5 := 2097151 & (load3(b[13:]) >> 1)
	b6 := 2097151 & (load4(b[15:]) >> 6)
	b7 := 2097151 & (load3(b[18:]) >> 3)
	b8 := 2097151 & load3(b[21:])
	b9 := 2097151 & (load4(b[23:]) >> 5)
	b10 := 2097151 & (load3(b[26:]) >> 2)
	b11 := load4(b[28:]) >> 7

	s0 := a0 - b0
	s1 := a1 - b1
	s2 := a2 - b2
	s3 := a3 - b3
	s4 := a4 - b4
	s5 := a5 - b5
	s6 := a6 - b6
	s7 := a7 - b7
	s8 := a8 - b8
	s9 := a9 - b9
	s10 := a10 - b10
	s11 := a11 - b11
	s12 := int64(0)

	carry0 := (s0 + (1 << 20)) >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry2 := (s2 + (1 << 20)) >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry4 := (s4 + (1 << 20)) >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry6 := (s6 + (1 << 20)) >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry8 := (s8 + (1 << 20)) >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry10 := (s10 + (1 << 20)) >> 21
	s11 += carry10
	s10 -= carry10 << 21

	carry1 := (s1 + (1 << 20)) >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry3 := (s3 + (1 << 20)) >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry5 := (s5 + (1 << 20)) >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry7 := (s7 + (1 << 20)) >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry9 := (s9 + (1 << 20)) >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry11 := (s11 + (1 << 20)) >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901
	s12 = 0

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21
	carry11 = s11 >> 21
	s12 += carry11
	s11 -= carry11 << 21

	s0 += s12 * 666643
	s1 += s12 * 470296
	s2 += s12 * 654183
	s3 -= s12 * 997805
	s4 += s12 * 136657
	s5 -= s12 * 683901

	carry0 = s0 >> 21
	s1 += carry0
	s0 -= carry0 << 21
	carry1 = s1 >> 21
	s2 += carry1
	s1 -= carry1 << 21
	carry2 = s2 >> 21
	s3 += carry2
	s2 -= carry2 << 21
	carry3 = s3 >> 21
	s4 += carry3
	s3 -= carry3 << 21
	carry4 = s4 >> 21
	s5 += carry4
	s4 -= carry4 << 21
	carry5 = s5 >> 21
	s6 += carry5
	s5 -= carry5 << 21
	carry6 = s6 >> 21
	s7 += carry6
	s6 -= carry6 << 21
	carry7 = s7 >> 21
	s8 += carry7
	s7 -= carry7 << 21
	carry8 = s8 >> 21
	s9 += carry8
	s8 -= carry8 << 21
	carry9 = s9 >> 21
	s10 += carry9
	s9 -= carry9 << 21
	carry10 = s10 >> 21
	s11 += carry10
	s10 -= carry10 << 21

	s[0] = byte(s0 >> 0)
	s[1] = byte(s0 >> 8)
	s[2] = byte((s0 >> 16) | (s1 << 5))
	s[3] = byte(s1 >> 3)
	s[4] = byte(s1 >> 11)
	s[5] = byte((s1 >> 19) | (s2 << 2))
	s[6] = byte(s2 >> 6)
	s[7] = byte((s2 >> 14) | (s3 << 7))
	s[8] = byte(s3 >> 1)
	s[9] = byte(s3 >> 9)
	s[10] = byte((s3 >> 17) | (s4 << 4))
	s[11] = byte(s4 >> 4)
	s[12] = byte(s4 >> 12)
	s[13] = byte((s4 >> 20) | (s5 << 1))
	s[14] = byte(s5 >> 7)
	s[15] = byte((s5 >> 15) | (s6 << 6))
	s[16] = byte(s6 >> 2)
	s[17] = byte(s6 >> 10)
	s[18] = byte((s6 >> 18) | (s7 << 3))
	s[19] = byte(s7 >> 5)
	s[20] = byte(s7 >> 13)
	s[21] = byte(s8 >> 0)
	s[22] = byte(s8 >> 8)
	s[23] = byte((s8 >> 16) | (s9 << 5))
	s[24] = byte(s9 >> 3)
	s[25] = byte(s9 >> 11)
	s[26] = byte((s9 >> 19) | (s10 << 2))
	s[27] = byte(s10 >> 6)
	s[28] = byte((s10 >> 14) | (s11 << 7))
	s[29] = byte(s11 >> 1)
	s[30] = byte(s11 >> 9)
	s[31] = byte(s11 >> 17)

	return s[:]
}
