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
	// Date is the same as ekatime.Date but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/date_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Date struct {
		ekatime_orig.Date
	}
)

var (
	_ types.ValueAppender = (*Date)(nil)
	_ types.ValueScanner  = (*Date)(nil)
)

func (dd *Date) AppendValue(b []byte, flags int) ([]byte, error) {

	if dd.Date == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = dd.Date.AppendTo(b, '-')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (dd *Date) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if dd != nil {
			dd.Date = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return dd.ParseFrom(b)
}

// WrapDate returns an Date object as modified ekatime.Date object for being able
// to use it with go-pg.
//
// See also: WrapDateUnsafe().
func WrapDate(dd ekatime_orig.Date) Date {
	return Date{dd}
}

// WrapDateUnsafe uses some Golang internal features and provides you a zero cost
// pointer cast conversion between standard ekatime.Date type and modified Date
// type that supports go-pg.
//
// See also: WrapDate().
func WrapDateUnsafe(dd *ekatime_orig.Date) *Date {
	return (*Date)(unsafe.Pointer(dd)) // it's ok, at least at the Go 1.13 - Go 1.15
}
