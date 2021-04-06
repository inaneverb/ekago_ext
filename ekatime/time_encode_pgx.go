package ekatime

import (
	"encoding/binary"
	"fmt"

	ekatime_orig "github.com/qioalice/ekago/v2/ekatime"

	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
)

const (
	microsecondsPerSecond = 1000000
	microsecondsPerMinute = 60 * microsecondsPerSecond
	microsecondsPerHour   = 60 * microsecondsPerMinute
)

// DecodeBinary
//
// From https://github.com/jackc/pgtype@v1.7.0/pgtype.go :
//
//     DecodeBinary decodes src into BinaryDecoder. If src is nil then the
//     original SQL value is NULL. BinaryDecoder takes ownership of src. The
//     caller MUST not use it again.
//
func (t *Time) DecodeBinary(_ *pgtype.ConnInfo, src []byte) error {

	if src == nil {
		*t = 0
		return nil
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for time: %v", len(src))
	}

	microsecondsSinceDayStart := int64(binary.BigEndian.Uint64(src))

	hh := microsecondsSinceDayStart / microsecondsPerHour
	microsecondsSinceDayStart -= hh * microsecondsPerHour

	mm := microsecondsSinceDayStart / microsecondsPerMinute
	microsecondsSinceDayStart -= mm * microsecondsPerMinute

	ss := microsecondsSinceDayStart / microsecondsPerSecond

	*t = WrapTime(ekatime_orig.NewTime(Hour(hh), Minute(mm), Second(ss)))

	return nil
}

// EncodeBinary
//
// From https://github.com/jackc/pgtype@v1.7.0/pgtype.go :
//
//     EncodeBinary should append the binary format of self to buf. If self is the
//     SQL value NULL then append nothing and return (nil, nil). The caller of
//     EncodeBinary is responsible for writing the correct NULL value or the
//     length of the data written.
//
func (t Time) EncodeBinary(_ *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {

	if t == 0 {
		return nil, nil
	}

	microsecondsSinceDayStart :=
		int64(t.Second()) * microsecondsPerSecond +
		int64(t.Minute()) * microsecondsPerMinute +
		int64(t.Hour()) * microsecondsPerHour

	return pgio.AppendInt64(buf, microsecondsSinceDayStart), nil
}
