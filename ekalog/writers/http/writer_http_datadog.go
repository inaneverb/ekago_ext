package ekalog_writer_http

import (
	"github.com/valyala/fasthttp"
)

//noinspection GoSnakeCaseUsage
const (
	DATADOG_ADDR_US = "https://http-intake.logs.datadoghq.com/v1/input"
	DATADOG_ADDR_EU = "https://http-intake.logs.datadoghq.eu/v1/input"
)

// UseProviderDataDog setups CI_WriterHttp for DataDog log service provider
// ( https://www.datadoghq.com/ ).
//
// You MUST specify 'addr' as desired DataDog's HTTP addr (you may use predefined
// constants DATADOG_ADDR_US, DATADOG_ADDR_EU) or use your own and DataDog
// service's token as 'token'.
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) UseProviderDataDog(addr, token string) *CI_WriterHttp {

	cb1 := func(req *fasthttp.Request) {
		req.SetRequestURI(addr)
		req.Header.SetContentType("application/json")
		req.Header.Set("DD-API-KEY", token)
	}

	return dw.AddBeforeAfterBetweenS("[],").UseProviderManual(cb1)
}
