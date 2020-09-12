// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"unsafe"

	ekatime_orig "github.com/qioalice/ekago/v2/ekatime"
	"github.com/qioalice/ekago_ext/v2/internal/ekasql"

	"github.com/go-pg/pg/v10/types"
)

type (
	// Timestamp is the same as ekatime.Timestamp but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/timestamp_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Timestamp struct {
		ekatime_orig.Timestamp
	}
)

var (
	_ types.ValueAppender = (*Timestamp)(nil)
	_ types.ValueScanner  = (*Timestamp)(nil)
)

func (ts *Timestamp) AppendValue(b []byte, flags int) ([]byte, error) {

	if ts.Timestamp == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = ts.Timestamp.AppendTo(b, '-', ':')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (ts *Timestamp) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if ts != nil {
			ts.Timestamp = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return ts.ParseFrom(b)
}

// WrapDate returns a Timestamp object as modified ekatime.Timestamp object
// for being able to use it with go-pg.
//
// See also: WrapDateUnsafe().
func WrapTimestamp(ts ekatime_orig.Timestamp) Timestamp {
	return Timestamp{ts}
}

// WrapDateUnsafe uses some Golang internal features and provides you a zero cost
// pointer cast conversion between standard ekatime.Timestamp type
// and modified Timestamp type that supports go-pg.
//
// See also: WrapDate().
func WrapTimestampUnsafe(ts *ekatime_orig.Timestamp) *Timestamp {
	return (*Timestamp)(unsafe.Pointer(ts)) // it's ok, at least at the Go 1.13 - Go 1.15
}
