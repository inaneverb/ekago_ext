package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/qioalice/ekago_ext/v2/internal/ekasql"

	"github.com/go-pg/pg/v10/types"
)

var (
	_ types.ValueAppender = (*Date)(nil)
	_ types.ValueScanner  = (*Date)(nil)
)

func (dd Date) AppendValue(b []byte, flags int) ([]byte, error) {

	if dd == 0 {
		return ekasql.NULL_AS_BYTES_SLICE, nil
	}

	if flags == 1 {
		b = append(b, '\'')
	}
	b = dd.AppendTo(b, '-')
	if flags == 1 {
		b = append(b, '\'')
	}

	return b, nil
}

func (dd *Date) ScanValue(rd types.Reader, n int) error {

	if n <= 0 {
		if dd != nil {
			*dd = 0
		}
		return nil
	}

	b, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	return dd.ParseFrom(b)
}

func (dd Date) AppendTo(b []byte, separator byte) []byte {
	return dd.ToOrig().AppendTo(b, separator)
}

func (dd *Date) ParseFrom(b []byte) error {
	return dd.ToOrigPtr().ParseFrom(b)
}

func (dd Date) String() string {
	return dd.ToOrig().String()
}

func (dd *Date) MarshalJSON() ([]byte, error) {
	return dd.ToOrigPtr().MarshalJSON()
}

func (dd *Date) UnmarshalJSON(b []byte) error {
	return dd.ToOrigPtr().UnmarshalJSON(b)
}
