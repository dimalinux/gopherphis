package polyseed

import (
	"time"
)

const (

	// UnixEpochDelta is the earliest time, as a delta in seconds from the Unix
	// Epoch, that can be represented by polyseed's time encoding: 1st November
	// 2021 12:00 UTC
	UnixEpochDelta = 1635768000

	// TimeStep is the resolution of a polyseed date in seconds
	TimeStep = 2629746 // 30.436875 days = 1/12 of the Gregorian year

	// DateBits the number of bits used to represent a polyseed date
	DateBits = 10

	// DateMask is a bit mask for the polyseed date bits
	DateMask = (1 << DateBits) - 1
)

// birthdayEncode converts epoch time into polyseed's lower granularity version of time
// that consumes at most 10 bits when encoded.
func birthdayEncode(epochTime int64) uint {
	if epochTime < UnixEpochDelta {
		// Just return the beginning of polytime. We were given non-nonsensical
		// input, as polytime didn't exist before UnixEpochDelta.
		return 0
	}

	polyTime := uint(epochTime-UnixEpochDelta) / TimeStep

	// If the date is consuming more than 10 bits (i.e. after February 2107), we
	// received non-nonsensical input, so just return the beginning of poly
	// time. The user's balance scanning will just be slower than it otherwise
	// would have been.
	if polyTime&DateMask != polyTime {
		polyTime = 0
	}

	return polyTime
}

func birthdayDecode(birthday uint) int64 {
	return UnixEpochDelta + int64(birthday&DateMask)*TimeStep
}

var epochTimeNow = func() int64 {
	return time.Now().Unix()
}

func birthdayNow() uint {
	const dayInSeconds = 60 * 60 * 24 // 1 day in seconds
	nowEpoch := epochTimeNow()
	polyTime := birthdayEncode(nowEpoch)

	// Each poly unit of time is a little over 30 days. As an extra safety for
	// computers that have the incorrect time or time zone, we subtract off a
	// single poly time unit if the poly time is within 24 hours of the actual
	// time.
	polyTimeInEpoch := birthdayDecode(polyTime)
	if nowEpoch-polyTimeInEpoch < dayInSeconds {
		polyTime--
	}

	return polyTime
}
