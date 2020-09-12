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
	// Time is the same as ekatime.Time but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/time_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Time struct {
		ekatime_orig.Time
	}
)

var (
	_ types.ValueAppender = (*Time)(nil)
	_ types.ValueScanner  = (*Time)(nil)
)

func (t *Time) AppendValue(b []byte, flags int) ([]byte, error) {

	if t.Time == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = t.Time.AppendTo(b, ':')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (t *Time) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if t != nil {
			t.Time = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return t.ParseFrom(b)
}

// WrapDate returns a Time object as modified ekatime.Time object for being able
// to use it with go-pg.
//
// See also: WrapDateUnsafe().
func WrapTime(ts ekatime_orig.Time) Time {
	return Time{ts}
}

// WrapDateUnsafe uses some Golang internal features and provides you a zero cost
// pointer cast conversion between standard ekatime.Time type and modified Time
// type that supports go-pg.
//
// See also: WrapDate().
func WrapTimeUnsafe(ts *ekatime_orig.Time) *Time {
	return (*Time)(unsafe.Pointer(ts)) // it's ok, at least at the Go 1.13 - Go 1.15
}
