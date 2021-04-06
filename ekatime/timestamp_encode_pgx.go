package ekatime

import (
	"encoding/binary"
	"fmt"

	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
)

const (
	microsecFromUnixEpochToY2K        = 946684800 * 1000000
	negativeInfinityMicrosecondOffset = -9223372036854775808
	infinityMicrosecondOffset         = 9223372036854775807
)

// DecodeBinary
//
// From https://github.com/jackc/pgtype@v1.7.0/pgtype.go :
//
//     DecodeBinary decodes src into BinaryDecoder. If src is nil then the
//     original SQL value is NULL. BinaryDecoder takes ownership of src. The
//     caller MUST not use it again.
//
func (ts *Timestamp) DecodeBinary(_ *pgtype.ConnInfo, src []byte) error {

	// Original code from:
	// https://github.com/jackc/pgtype/blob/4a3a424/timestamp.go (v1.7.0)

	if src == nil {
		*ts = 0
		return nil
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for timestamp: %v", len(src))
	}

	microsecSinceY2K := int64(binary.BigEndian.Uint64(src))

	switch microsecSinceY2K {
	case infinityMicrosecondOffset, negativeInfinityMicrosecondOffset:
		*ts = Timestamp(microsecSinceY2K)
	default:
		microsecSinceUnixEpoch := microsecFromUnixEpochToY2K + microsecSinceY2K
		*ts = Timestamp(microsecSinceUnixEpoch/1000000)
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
func (ts Timestamp) EncodeBinary(_ *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {

	// Original code from:
	// https://github.com/jackc/pgtype/blob/4a3a424/timestamp.go (v1.7.0)

	if ts == 0 {
		return nil, nil
	}

	var microsecSinceY2K int64

	switch ts {
	case infinityMicrosecondOffset, negativeInfinityMicrosecondOffset:
		microsecSinceY2K = ts.I64()
	default:
		microsecSinceUnixEpoch := ts.I64() * 1000000
		microsecSinceY2K = microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	}

	return pgio.AppendInt64(buf, microsecSinceY2K), nil
}
