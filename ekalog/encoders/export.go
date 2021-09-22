// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_encoders

import (
	"github.com/qioalice/ekago/v3/ekalog"

	"github.com/qioalice/ekago_ext/v3/ekalog/encoders/datadog"
)

func NewDatadogJsonEncoder() *ekalog.CI_JSONEncoder {
	return ekalog_encoder_datadog.NewJsonEncoder()
}
