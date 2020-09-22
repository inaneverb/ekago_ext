// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	ekatime_orig "github.com/qioalice/ekago/v2/ekatime"
)

type (
	// Time is the same as ekatime.Time but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/time_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Time ekatime_orig.Time
)

// WrapDate returns a Time object as modified ekatime.Time object for being able
// to use it with go-pg.
//
// See also: WrapDateUnsafe().
func WrapTime(t ekatime_orig.Time) Time {
	return Time(t)
}

// WrapTimePtr returns a Time object by ptr as modified ekatime.Time object
// for being able to use it with go-pg.
//
// See also: WrapTime().
func WrapTimePtr(t *ekatime_orig.Time) *Time {
	return (*Time)(t)
}

func (t Time) ToOrig() ekatime_orig.Time {
	return ekatime_orig.Time(t)
}

func (t *Time) ToOrigPtr() *ekatime_orig.Time {
	return (*ekatime_orig.Time)(t)
}

func (t Time) Hour() Hour {
	return t.ToOrig().Hour()
}

func (t Time) Minute() Minute {
	return t.ToOrig().Minute()
}

func (t Time) Second() Second {
	return t.ToOrig().Second()
}

func (t Time) Split() (h Hour, m Minute, s Second) {
	return t.ToOrig().Split()
}

func NewTime(h Hour, m Minute, s Second) Time {
	return WrapTime(ekatime_orig.NewTime(h, m, s))
}

func (t Time) WithDate(y Year, m Month, d Day) Timestamp {
	return WrapTimestamp(t.ToOrig().WithDate(y, m, d))
}
