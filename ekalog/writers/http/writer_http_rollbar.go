// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_writer_http

import (
	"bytes"
	"io"

	"github.com/qioalice/ekago/v3/ekastr"

	"github.com/valyala/fasthttp"
)

//noinspection GoSnakeCaseUsage
const (
	ROLLBAR_ADDR = "https://api.rollbar.com/api/1/items/"
)

func (dw *CI_WriterHttp) UseProviderRollbar(token string) *CI_WriterHttp {

	var (
		JsonBodyStart = ekastr.S2B(`"access_token":"` + token + `",data:`)
		JsonBodyEnd   = ekastr.S2B(`}`)
	)

	cb1 := func(req *fasthttp.Request) {
		req.SetRequestURI(ROLLBAR_ADDR)
		req.Header.SetContentType("application/json")
	}

	cb2 := func(oldBody io.Reader) io.Reader {
		r := bytes.NewReader
		return io.MultiReader(r(JsonBodyStart), oldBody, r(JsonBodyEnd))
	}

	return dw.AddBeforeAfterBetweenS("[],").UseProviderManual(cb1, cb2)
}
