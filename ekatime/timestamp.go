// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"time"

	ekatime_orig "github.com/qioalice/ekago/v3/ekatime"
)

type (
	// Timestamp is the same as ekatime.Timestamp but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/timestamp_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	Timestamp OriginalTimestamp
)

// WrapTimestamp returns a Timestamp object as modified ekatime.Timestamp object
// for being able to use it with go-pg.
//
// See also: WrapTimestampPtr().
func WrapTimestamp(ts OriginalTimestamp) Timestamp {
	return Timestamp(ts)
}

// WrapTimestampPtr returns a Timestamp object by ptr as modified ekatime.Time object
// for being able to use it with go-pg.
//
// See also: WrapTimestamp().
func WrapTimestampPtr(ts *OriginalTimestamp) *Timestamp {
	return (*Timestamp)(ts)
}

func (ts Timestamp) ToOrig() OriginalTimestamp {
	return OriginalTimestamp(ts)
}

func (ts *Timestamp) ToOrigPtr() *OriginalTimestamp {
	return (*OriginalTimestamp)(ts)
}

func (ts Timestamp) I64() int64 {
	return ts.ToOrig().I64()
}

func (ts Timestamp) Std() time.Time {
	return ts.ToOrig().Std()
}

func (ts Timestamp) Date() Date {
	return WrapDate(ts.ToOrig().Date())
}

func (ts Timestamp) Time() Time {
	return WrapTime(ts.ToOrig().Time())
}

func (ts Timestamp) Year() Year {
	return ts.ToOrig().Year()
}

func (ts Timestamp) Month() Month {
	return ts.ToOrig().Month()
}

func (ts Timestamp) Day() Day {
	return ts.ToOrig().Day()
}

func (ts Timestamp) Hour() Hour {
	return ts.ToOrig().Hour()
}

func (ts Timestamp) Minute() Minute {
	return ts.ToOrig().Minute()
}

func (ts Timestamp) Second() Second {
	return ts.ToOrig().Second()
}

func (ts Timestamp) Split() (d Date, t Time) {
	return WrapDate(ts.ToOrig().Date()), WrapTime(ts.ToOrig().Time())
}

func NewTimestampNow() Timestamp {
	return WrapTimestamp(ekatime_orig.NewTimestampNow())
}

func NewTimestamp(y Year, m Month, d Day, hh Hour, mm Minute, ss Second) Timestamp {
	return WrapTimestamp(ekatime_orig.NewTimestamp(y, m, d, hh, mm, ss))
}

func NewTimestampFromStd(t time.Time) Timestamp {
	return WrapTimestamp(ekatime_orig.NewTimestampFromStd(t))
}

func (ts Timestamp) BeginningOfDay() Timestamp {
	return WrapTimestamp(ts.ToOrig().BeginningOfDay())
}

func (ts Timestamp) EndOfDay() Timestamp {
	return WrapTimestamp(ts.ToOrig().EndOfDay())
}

func (ts Timestamp) BeginningAndEndOfDay() TimestampPair {
	return WrapTimestampPair(ts.ToOrig().BeginningAndEndOfDay())
}

func (ts Timestamp) BeginningOfMonth() Timestamp {
	return WrapTimestamp(ts.ToOrig().BeginningOfMonth())
}

func (ts Timestamp) EndOfMonth() Timestamp {
	return WrapTimestamp(ts.ToOrig().EndOfMonth())
}

func (ts Timestamp) BeginningAndEndOfMonth() TimestampPair {
	return WrapTimestampPair(ts.ToOrig().BeginningAndEndOfMonth())
}

func (ts Timestamp) BeginningOfYear() Timestamp {
	return WrapTimestamp(ts.ToOrig().BeginningOfYear())
}

func (ts Timestamp) EndOfYear() Timestamp {
	return WrapTimestamp(ts.ToOrig().EndOfYear())
}

func (ts Timestamp) BeginningAndEndOfYear() TimestampPair {
	return WrapTimestampPair(ts.ToOrig().BeginningAndEndOfYear())
}

func BeginningOfYear(y Year) Timestamp {
	return WrapTimestamp(ekatime_orig.BeginningOfYear(y))
}

func EndOfYear(y Year) Timestamp {
	return WrapTimestamp(ekatime_orig.EndOfYear(y))
}

func BeginningAndEndOfYear(y Year) TimestampPair {
	return WrapTimestampPair(ekatime_orig.BeginningAndEndOfYear(y))
}

func BeginningOfMonth(y Year, m Month) Timestamp {
	return WrapTimestamp(ekatime_orig.BeginningOfMonth(y, m))
}

func EndOfMonth(y Year, m Month) Timestamp {
	return WrapTimestamp(ekatime_orig.EndOfMonth(y, m))
}

func BeginningAndEndOfMonth(y Year, m Month) TimestampPair {
	return WrapTimestampPair(ekatime_orig.BeginningAndEndOfMonth(y, m))
}
