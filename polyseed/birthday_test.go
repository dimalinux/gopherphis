package polyseed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Tests converting from the current time, to polyseed time, and back. Since the
// granularity of polyseed time is only around 30 days, the reversed time will
// almost certainly be in the past.
func TestBirthdayNow(t *testing.T) {
	now := time.Now()
	polyNow := birthdayNow()
	nowReversed := time.Unix(birthdayDecode(polyNow), 0)
	require.LessOrEqual(t, nowReversed, now)
	require.LessOrEqual(t, now.Sub(nowReversed), TimeStep*time.Second)
}

func Test_birthdayExactlyOnPolyseedStartOfTime(t *testing.T) {
	require.Zero(t, birthdayEncode(UnixEpochDelta))
	require.Equal(t, birthdayDecode(0), time.Unix(UnixEpochDelta, 0).Unix())
}

// Tests that any value before the start of polyseed time or consuming more than 10 bits
// just return the polyseed start of time.
func Test_birthdayOutOfRange(t *testing.T) {
	// 20 years before the Unix epoch
	beforeEpoch := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	require.Zero(t, birthdayEncode(beforeEpoch))

	// 1 second before the polyseed start of time
	require.Zero(t, birthdayEncode(time.Unix(UnixEpochDelta-1, 0).Unix()))

	// A full polyseed time unit in epoch seconds past the max polyseed time of
	// 1023. Encoded polyseed time will be 1024, exceeding 10 bits, and treated
	// as nonsensical, so the encoded birthday is just the start of poly time.
	require.Zero(t, birthdayEncode((DateMask+1)*TimeStep+UnixEpochDelta))

}
