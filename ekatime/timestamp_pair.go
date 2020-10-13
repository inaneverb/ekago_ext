// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
import (
	"unsafe"
)

type (
	// TimestampPair is the same as ekatime.TimestampPair
	// but with supporting go-pg (v10).
	//
	// Read more:
	// https://github.com/qioalice/ekago/ekatime/timestamp_encode.go ,
	// https://github.com/go-pg/pg ,
	// https://github.com/go-pg/pg/blob/v10/example_custom_test.go .
	TimestampPair [2]Timestamp
)

func WrapTimestampPair(tsp OriginalTimestampPair) TimestampPair {
	return *(*TimestampPair)(unsafe.Pointer(&tsp))
}

func (tsp TimestampPair) ToOrig() OriginalTimestampPair {
	return *(*OriginalTimestampPair)(unsafe.Pointer(&tsp))
}

func (tsp TimestampPair) Split() (t1, t2 Timestamp) {
	return tsp[0], tsp[1]
}

func (tsp TimestampPair) I64() (int64, int64) {
	return int64(tsp[0]), int64(tsp[1])
}
