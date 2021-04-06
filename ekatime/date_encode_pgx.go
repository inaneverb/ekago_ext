package ekatime

import (
	"encoding/binary"
	"fmt"
	"time"

	ekatime_orig "github.com/qioalice/ekago/v2/ekatime"

	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
)

//goland:noinspection GoSnakeCaseUsage
const (
	negativeInfinityDayOffset = -2147483648
	infinityDayOffset         = 2147483647

	// _TIMESTAMP_2000_Jan_01 is a Unix timestamp of 01 Jan 2000 in UTC
	_TIMESTAMP_2000_Jan_01 Timestamp = 946684800

	// Add constants from original ekatime package
	// [24..31] bits are unused and reserved for internal purposes.
	_DATE_INF = Date(1) << 26
	_DATE_IS_NEG_INF = Date(1) << 27
)

// DecodeBinary
//
// From https://github.com/jackc/pgtype@v1.7.0/pgtype.go :
//
//     DecodeBinary decodes src into BinaryDecoder. If src is nil then the
//     original SQL value is NULL. BinaryDecoder takes ownership of src. The
//     caller MUST not use it again.
//
func (dd *Date) DecodeBinary(_ *pgtype.ConnInfo, src []byte) error {

	// Original code from:
	// https://github.com/jackc/pgtype/blob/4a3a424/date.go (v1.7.0)

	if src == nil {
		*dd = 0
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for date: %v", len(src))
	}

	dayOffset := int32(binary.BigEndian.Uint32(src))

	switch dayOffset {
	case infinityDayOffset:
		*dd = 0 | _DATE_INF
	case negativeInfinityDayOffset:
		*dd = 0 | _DATE_INF | _DATE_IS_NEG_INF
	default:
		t := time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC)
		*dd = WrapDate(ekatime_orig.UnixFromStd(t).Date())
	}

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
func (dd Date) EncodeBinary(_ *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {

	// Original code from:
	// https://github.com/jackc/pgtype/blob/4a3a424/date.go (v1.7.0)

	if dd == 0 {
		return nil, nil
	}

	var daysSinceDateEpoch int32

	switch {
	case dd & _DATE_INF != 0 && dd & _DATE_IS_NEG_INF != 0:
		daysSinceDateEpoch = negativeInfinityDayOffset
	case dd & _DATE_INF != 0:
		daysSinceDateEpoch = infinityDayOffset
	default:
		tUnix := dd.WithTime(0, 0, 0)
		secSinceDateEpoch := tUnix - _TIMESTAMP_2000_Jan_01
		daysSinceDateEpoch = int32(secSinceDateEpoch / 86400)
	}

	return pgio.AppendInt32(buf, daysSinceDateEpoch), nil
}
