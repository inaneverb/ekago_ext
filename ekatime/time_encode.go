// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/qioalice/ekago_ext/v3/internal/ekasql"

	"github.com/go-pg/pg/v10/types"
)

var (
	_ types.ValueAppender = (*Time)(nil)
	_ types.ValueScanner  = (*Time)(nil)
)

func (t Time) AppendValue(b []byte, flags int) ([]byte, error) {

	if t == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = t.AppendTo(b, ':')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (t *Time) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if t != nil {
			*t = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return t.ParseFrom(b)
}

func (t Time) AppendTo(b []byte, separator byte) []byte {
	return t.ToOrig().AppendTo(b, separator)
}

func (t *Time) ParseFrom(b []byte) error {
	return t.ToOrigPtr().ParseFrom(b)
}

func (t Time) String() string {
	return t.ToOrig().String()
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return t.ToOrigPtr().MarshalJSON()
}

func (t *Time) UnmarshalJSON(b []byte) error {
	return t.ToOrigPtr().UnmarshalJSON(b)
}
