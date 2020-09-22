// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/qioalice/ekago_ext/v2/internal/ekasql"

	"github.com/go-pg/pg/v10/types"
)

var (
	_ types.ValueAppender = (*Timestamp)(nil)
	_ types.ValueScanner  = (*Timestamp)(nil)
)

func (ts Timestamp) AppendValue(b []byte, flags int) ([]byte, error) {

	if ts == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = ts.AppendTo(b, '-', ':')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (ts *Timestamp) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if ts != nil {
			*ts = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return ts.ParseFrom(b)
}

func (ts Timestamp) AppendTo(b []byte, separatorDate, separatorTime byte) []byte {
	return ts.ToOrig().AppendTo(b, separatorDate, separatorTime)
}

func (ts *Timestamp) ParseFrom(b []byte) error {
	return ts.ToOrigPtr().ParseFrom(b)
}

func (ts Timestamp) String() string {
	return ts.ToOrig().String()
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	return ts.ToOrigPtr().MarshalJSON()
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	return ts.ToOrigPtr().UnmarshalJSON(b)
}
