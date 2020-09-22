// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

func (ts Timestamp) TillNextMinute() time.Duration {
	return ts.ToOrig().TillNextMinute()
}

func (ts Timestamp) TillNextHour() time.Duration {
	return ts.ToOrig().TillNextHour()
}

func (ts Timestamp) TillNext12h() time.Duration {
	return ts.ToOrig().TillNext12h()
}

func (ts Timestamp) TillNextNoon() time.Duration {
	return ts.ToOrig().TillNextNoon()
}

func (ts Timestamp) TillNextMidnight() time.Duration {
	return ts.ToOrig().TillNextMidnight()
}

func (ts Timestamp) TillNextDay() time.Duration {
	return ts.ToOrig().TillNextDay()
}

func (ts Timestamp) TillNextMonth() time.Duration {
	return ts.ToOrig().TillNextMonth()
}

func (ts Timestamp) TillNextYear() time.Duration {
	return ts.ToOrig().TillNextYear()
}
