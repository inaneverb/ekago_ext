// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	ekatime_orig "github.com/qioalice/ekago/v2/ekatime"
)

type (
	// https://github.com/qioalice/ekago/ekatime/date.go

	Year  = ekatime_orig.Year
	Month = ekatime_orig.Month
	Day   = ekatime_orig.Day

	// https://github.com/qioalice/ekago/ekatime/time.go

	Hour   = ekatime_orig.Hour
	Minute = ekatime_orig.Minute
	Second = ekatime_orig.Second

	// https://github.com/qioalice/ekago/ekatime/weekday.go

	Weekday = ekatime_orig.Weekday

	// https://github.com/qioalice/ekago/ekatime/event.go

	Event = ekatime_orig.Event

	// https://github.com/qioalice/ekago/ekatime/today.go

	Today = ekatime_orig.Today

	// https://github.com/qioalice/ekago/ekatime/calendar.go

	Calendar = ekatime_orig.Calendar
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	// https://github.com/qioalice/ekago/ekatime/date.go

	MONTH_JANUARY   = ekatime_orig.MONTH_JANUARY
	MONTH_FEBRUARY  = ekatime_orig.MONTH_FEBRUARY
	MONTH_MARCH     = ekatime_orig.MONTH_MARCH
	MONTH_APRIL     = ekatime_orig.MONTH_APRIL
	MONTH_MAY       = ekatime_orig.MONTH_MAY
	MONTH_JUNE      = ekatime_orig.MONTH_JUNE
	MONTH_JULY      = ekatime_orig.MONTH_JULY
	MONTH_AUGUST    = ekatime_orig.MONTH_AUGUST
	MONTH_SEPTEMBER = ekatime_orig.MONTH_SEPTEMBER
	MONTH_OCTOBER   = ekatime_orig.MONTH_OCTOBER
	MONTH_NOVEMBER  = ekatime_orig.MONTH_NOVEMBER
	MONTH_DECEMBER  = ekatime_orig.MONTH_DECEMBER

	// https://github.com/qioalice/ekago/ekatime/weekday.go

	WEEKDAY_WEDNESDAY = ekatime_orig.WEEKDAY_WEDNESDAY
	WEEKDAY_THURSDAY  = ekatime_orig.WEEKDAY_THURSDAY
	WEEKDAY_FRIDAY    = ekatime_orig.WEEKDAY_FRIDAY
	WEEKDAY_SATURDAY  = ekatime_orig.WEEKDAY_SATURDAY
	WEEKDAY_SUNDAY    = ekatime_orig.WEEKDAY_SUNDAY
	WEEKDAY_MONDAY    = ekatime_orig.WEEKDAY_MONDAY
	WEEKDAY_TUESDAY   = ekatime_orig.WEEKDAY_TUESDAY

	// https://github.com/qioalice/ekago/ekatime/internal.go

	SECONDS_IN_MINUTE   = ekatime_orig.SECONDS_IN_MINUTE
	SECONDS_IN_HOUR     = ekatime_orig.SECONDS_IN_HOUR
	SECONDS_IN_12H      = ekatime_orig.SECONDS_IN_12H
	SECONDS_IN_DAY      = ekatime_orig.SECONDS_IN_DAY
	SECONDS_IN_WEEK     = ekatime_orig.SECONDS_IN_WEEK
	SECONDS_IN_365_YEAR = ekatime_orig.SECONDS_IN_365_YEAR
	SECONDS_IN_366_YEAR = ekatime_orig.SECONDS_IN_366_YEAR
)

//goland:noinspection GoUnusedGlobalVariable
var (
	// https://github.com/qioalice/ekago/ekatime/once_in.go
	OnceInMinute    = &ekatime_orig.OnceInMinute
	OnceIn10Minutes = &ekatime_orig.OnceIn10Minutes
	OnceIn15Minutes = &ekatime_orig.OnceIn15Minutes
	OnceIn30Minutes = &ekatime_orig.OnceIn30Minutes
	OnceInHour      = &ekatime_orig.OnceInHour
	OnceIn12Hours   = &ekatime_orig.OnceIn12Hours
	OnceInDay       = &ekatime_orig.OnceInDay
)

//goland:noinspection GoUnusedGlobalVariable
var (
	// https://github.com/qioalice/ekago/ekatime/date.go

	IsValidDate = ekatime_orig.IsValidDate
	IsLeap      = ekatime_orig.IsLeap

	// https://github.com/qioalice/ekago/ekatime/event.go

	NewEvent = ekatime_orig.NewEvent

	// https://github.com/qioalice/ekago/ekatime/till_next.go

	TillNextMinute   = ekatime_orig.TillNextMinute
	TillNextHour     = ekatime_orig.TillNextHour
	TillNextNoon     = ekatime_orig.TillNextNoon
	TillNextMidnight = ekatime_orig.TillNextMidnight
	TillNextDay      = ekatime_orig.TillNextDay
	TillNextMonth    = ekatime_orig.TillNextMonth
	TillNextYear     = ekatime_orig.TillNextYear

	// https://github.com/qioalice/ekago/ekatime/time.go

	IsValidTime = ekatime_orig.IsValidTime
)
