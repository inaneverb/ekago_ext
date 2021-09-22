// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_encoder_datadog

import (
	"time"

	"github.com/qioalice/ekago/v3/ekalog"
)

// NewJsonEncoder creates a Datadog encoder that is based on ekalog.CI_JSONEncoder.
// It has almost the same rules but with the following changes:
//
// 1. You can not set an indentation. It's always 0. No tabs, no new lines.
//    Even at the end of data buffer, that contains JSON encoded log entry.
//
// 2. Log entry's timestamp has a different name: "timestamp_real".
//
// 5. All log's fields, and all attached error's fields (each stack frame's fields)
//    is encoded as JSON key-value pair at the root.
//
// 6. Stacktrace (log's or attached error's) is encoded
//    as JSON array of strings to the JSON's root
//    (because DataDog does not supports arrays of objects).
//    Each string will represent stack frame (caller) in the following format:
//
//        "[<stack_index>]: <func_name_with_fullpath>(<short_filename>:<file_line>)"
//
// 7. Attached error's messages (each stack frame's message) are encoded
//    as JSON array of strings to the JSON's root.
//    (because DataDog does not supports arrays of objects).
//    Each string will represent stack frame's message in the following format:
//
//        "[<stack_index>]: <message>"
//
func NewJsonEncoder() *ekalog.CI_JSONEncoder {
	return new(ekalog.CI_JSONEncoder).
		SetOneDepthLevel(true).
		SetNameForField(ekalog.CI_JSON_ENCODER_FIELD_TIME, "timestamp_real").
		SetTimeFormatter(func(t time.Time) string {
			const ISO8601 = "2006-01-02T15:04:05.000-0700"
			return t.Format(ISO8601)
		})
}
