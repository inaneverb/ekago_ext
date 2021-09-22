// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	ekatime_orig "github.com/qioalice/ekago/v3/ekatime"
)

type (
	// Date is the same as ekatime.Date but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/date_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Date OriginalDate
)

// WrapDate returns an Date object as modified ekatime.Date object for being able
// to use it with go-pg.
//
// See also: WrapDatePtr().
func WrapDate(dd OriginalDate) Date {
	return Date(dd)
}

// WrapDatePtr returns a Date object by ptr as modified ekatime.Date object
// for being able to use it with go-pg.
//
// See also: WrapDate().
func WrapDatePtr(dd *OriginalDate) *Date {
	return (*Date)(dd)
}

func (dd Date) ToOrig() OriginalDate {
	return OriginalDate(dd)
}

func (dd *Date) ToOrigPtr() *OriginalDate {
	return (*OriginalDate)(dd)
}

func (dd Date) ToCmp() Date {
	return WrapDate(dd.ToOrig().ToCmp())
}

func (dd Date) Equal(other Date) bool {
	return dd.ToOrig().Equal(other.ToOrig())
}

func (dd Date) Year() Year {
	return dd.ToOrig().Year()
}

func (dd Date) Month() Month {
	return dd.ToOrig().Month()
}

func (dd Date) Day() Day {
	return dd.ToOrig().Day()
}

func (dd Date) DaysInMonth() Day {
	return dd.ToOrig().DaysInMonth()
}

func (dd Date) Weekday() Weekday {
	return dd.ToOrig().Weekday()
}

func (dd Date) Split() (y Year, m Month, d Day) {
	return dd.ToOrig().Split()
}

func NewDate(y Year, m Month, d Day) Date {
	return WrapDate(ekatime_orig.NewDate(y, m, d))
}

func NewDateFromDayOfYear(y Year, days ekatime_orig.Days) Date {
	return WrapDate(ekatime_orig.NewDateFromDayOfYear(y, days))
}

func (dd Date) WithTime(hh Hour, mm Minute, ss Second) Timestamp {
	return WrapTimestamp(dd.ToOrig().WithTime(hh, mm, ss))
}

func InMonth(y Year, m Month) Timestamp {
	return WrapTimestamp(ekatime_orig.InMonth(y, m))
}

func InYear(y Year) Timestamp {
	return WrapTimestamp(ekatime_orig.InYear(y))
}
