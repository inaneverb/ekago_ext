// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_writer_http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/ekastr"

	"github.com/valyala/fasthttp"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_WriterHttp is a type that implements an io.Writer - legacy Golang interface,
	// doing write encoded log's entry as []byte to the some HTTP service using
	// desired API.
	//
	// Features:
	// -----------
	//
	// 1. Async transport.
	//    When you calling Write() it just pushes encoded entry (as []byte)
	//    to the worker and does not blocks the routine.
	//    Spawn as many workers (at the initialization) as you want.
	//    See SetWorkersNum() method.
	//
	// 2. Thread-safe.
	//    You may call CI_WriterHttp to as many goroutines as you want.
	//    The CommonIntegrator under the hood just calls Write() for all
	//    writers they holds on.
	//
	// 3. Accumulated bulk deferred requests.
	//    Has an internal buffer of encoded log entries (as []byte)
	//    (see SetBufferCap() method), each worker pulling data from
	//    and putting to its own internal buffer (see SetWorkerBufferCap() method).
	//
	//    When worker buffer is full or when it's time to flush the accumulated data
	//    (the stopping of app or flush time has come), sends accumulated data
	//    to the HTTP service.
	//
	// 4. Control how encoded log entries will be combined before sending.
	//    Each worker aggregates a single encoded log entries with others
	//    to the log entries pack (contains many encoded log entries) before send it.
	//
	//    You may set HOW they will be combined:
	//
	//        - Need to put something to buffer before first encoded log entry added?
	//          (see AddBefore() method).
	//
	//        - Need to put something to buffer after last encoded log entry added?
	//          (see AddAfter() method).
	//
	//        - Need to put something to buffer between encoded log entries?
	//          I mean between all of them, excluding before first and after last.
	//          (see AddBetween() method).
	//
	// 5. Flushing data each N time intervals.
	//    Each worker flushes accumulated data on demand (the internal buffer is full)
	//    or at the timeout you may set (see SetWorkerAutoFlushDelay() method).
	//
	// 6. Slow down when service unavailable or network error.
	//    If request sending at the some worker is failed
	//    (it's only network or log service issue if you configured writer well),
	//    the usage of HTTP API service will slow down, logging what happens
	//    and why and saving (not dropping) your logs to the internal buffer.
	//
	//    They will send when connection will be restored.
	//    You may set the deferred log entries pack buffer capacity
	//    (pack, because each worker builds a pack of log entries, remember? see p.3)
	//    (see SetDeferredBufferCap() method).
	//
	// 7. Graceful shutdown.
	//    Of course, if you're familiar of ekadeath package. If you're not yet,
	//    it's time to: https://github.com/qioalice/ekago/ekadeath .
	//
	//    When you calling ekadeath.Die(), ekadeath.Exit() or writing a log
	//    with the level that marked as fatal, you won't lost aggregated logs!
	//
	//    The package will sends the rest of logs for the last time for you.
	//    And you do not need to do something for that.
	//
	//    Need more?
	//    No problem, RegisterGracefulShutdown() allows you to specify context,
	//    using which you may finally disable CI_WriterHttp
	//    and a sync.WaitGroup, using which you may be sure, that you get your control
	//    only when all accumulated logs are flushed.
	//    You may specify only one of them.
	//
	// 8. Auto-initialization:
	//    You do need to call methods like Start() or something like that.
	//    You may check whether you configuration is valid calling Ping() method,
	//    but it's not necessary.
	//
	//    Just call all configuration methods with the chaining style and pass
	//    CI_WriterHttp object to the CommonIntegrator's WriteTo() method
	//    (or to your own logging integrator) and there is!
	//    The CI_WriterHttp will be initialized at the first Write() call.
	//
	// 9. Configurable to use any service.
	//    It's a very customizable type using which you may stream your logs safely
	//    to the log aggregation services like:
	//         - DataDog: https://www.datadoghq.com/ : UseProviderDataDog(),
	//         - Rollbar: https://rollbar.com/       : UseProviderRollbar(),
	//         - GrayLog: https://www.graylog.org/   : UseProviderGrayLog(),
	//         - Sentry:  https://sentry.io/         : UseProviderSentry(),
	//         etc.
	//
	//     For those services, CI_WriterHttp has a methods (3rd column)
	//     that makes it easy to configure it for, but you may configure it manually.
	//
	//     If you want to manually set HTTP service use UseProviderManual() method,
	//     providing a fasthttp's Request initializer for your service.
	//
	// 10. Fast.
	//     Uses fasthttp ( https://github.com/valyala/fasthttp ) under the hood,
	//     as http client. Pools, reusing, caching, optimizations. All you need.
	//
	// --------
	//
	// WARNING!
	// DO NOT CALL Write() or Ping() METHODS UNTIL YOU FINISH ALL PREPARATIONS!
	// DO NOT PASS WRITER TO THE CommonIntegrator's WriteTo() METHOD UNTIL
	// YOU FINISH ALL PREPARATIONS!
	// IF YOU DO, THE CHANGES WILL NOT BE SAVED!
	//
	// WARNING! PANIC CAUTION!
	// YOU MUST SET THE LOG SERVICE YOU WANT TO WRITE LOG ENTRIES TO.
	// CHOOSE A PREDEFINED OR USE YOUR OWN.
	// IF YOU DO NOT DO THAT, THE INITIALIZATION WILL PANIC!
	//
	CI_WriterHttp struct {

		// Has getter or/and setter

		providerInitializer  func(req *fasthttp.Request)
		providerBodyPreparer func(oldBody io.Reader) (newBody io.Reader)

		entriesBufferLen         uint32
		deferredEntriesBufferLen *uint32

		workerNum              uint16
		workerEntriesBufferLen uint16
		workerFlushDelay       time.Duration

		dataBefore  []byte
		dataAfter   []byte
		dataBetween []byte

		workerFlushDeferredPerIter uint16

		// Internal parts

		casInitStatus int32
		slowInit      sync.Mutex

		beenPinged bool

		ctx        context.Context
		cancelFunc context.CancelFunc

		workersWg  sync.WaitGroup
		externalWg *sync.WaitGroup

		workerTickers []*time.Ticker

		// This channel will never be closed.
		entries             chan []byte
		entriesPackDeferred chan *bytes.Buffer

		entriesCompletelyLostCounter uint64

		c fasthttp.Client
	}
)

var (
	ErrWriterIsNil      = fmt.Errorf("CI_WriterHttp: writer is nil (not initialized)")
	ErrWriterDisabled   = fmt.Errorf("CI_WriterHttp: writer is disabled (stopped)")
	ErrWriterBufferFull = fmt.Errorf("CI_WriterHttp: writer's buffer is full")
)

// UseProviderManual is a log service provider manual configurator.
// You MUST specify a callback that will set-up an HTTP request for desired provider.
//
// A 2nd argument allows you to modify generated HTTP request's body
// before it will be sent. Your callback (if provided) must return a new body,
// that will be used instead of old one.
//
// WARNING!
// You MUST NOT save and reuse old body that you receiver in your 2nd callback!
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) UseProviderManual(

	cb func(req *fasthttp.Request),
	bodyPreparer ...func(reader io.Reader) io.Reader,

) *CI_WriterHttp {

	return dw.configure(func(dw *CI_WriterHttp) {
		dw.providerInitializer = cb
		if len(bodyPreparer) > 0 && bodyPreparer[0] != nil {
			dw.providerBodyPreparer = bodyPreparer[0]
		}
	})
}

// RegisterGracefulShutdown allows you to pass context.Context and sync.WaitGroup,
// that will be used to provide you graceful shutdown, meaning:
//
// 1. Context.
//    Specify, when running CI_WriterHttp must be disabled.
//
// 2. sync.WaitGroup.
//    If specified, one second before all workers are started,
//    your waitgroup's counter will be increased, and it will be decreased,
//    when all of them are stopped, guaranteeing to you,
//    when disabling is requested, flushing all accumulated logs is important.
//
// Read p.1, p.8 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// You may pass only context or only sync.WaitGroup. It's OK.
func (dw *CI_WriterHttp) RegisterGracefulShutdown(ctx context.Context, wg *sync.WaitGroup) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.ctx = ctx
		dw.externalWg = wg
	})
}

// SetBufferCap sets a limit of internal pool of encoded []byte entries,
// to which Write() method places them, and where they are extracted from later
// for being processed and sent.
//
// If this cap is reached, Write() will be IGNORED all next entries,
// until old ones are processed.
//
// So, you need to specify this buffer as big as it must be never reached,
// to keep your logs safe even at the load,
// but not as big to spend gigabytes of RAM for them.
// Keep in mind, that buffer stores only pointers to []byte
// (it's up to 24 bytes per item), but the data of []byte also takes a RAM.
//
// Read p.1, p.3 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// Allowed range: [256..1'048'576] (2**8..2**20).
// Do not overwrite default value if you're not understand how much RAM
// your log entries consume per item, and what RAM consumption will be
// at the upper bound.
// Default: 4096.
func (dw *CI_WriterHttp) SetBufferCap(cap uint32) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		if cap >= uint32(1<<8) && cap <= uint32(1<<20) {
			dw.entriesBufferLen = cap
		}
	})
}

// SetWorkerBufferCap sets a limit of each worker's internal pool
// of encoded []byte entries, when that accumulated set will be sent to your provider.
//
// Less entries may be sent
// (if timeout that you may set by SetWorkerAutoFlushDelay() is reached),
// but this value tells a maximum.
//
// Read p.1, p.3 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// Allowed range: [1..16384].
// Default: 32.
//
// Hint.
// It's useful when your service may accept a bulk requests.
// Do less requests but much loaded.
// But if your provider doesn't accept bulk requests, you may need to set this to 1.
func (dw *CI_WriterHttp) SetWorkerBufferCap(cap uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		if cap >= 1 && cap <= 16384 {
			dw.workerEntriesBufferLen = cap
		}
	})
}

// SetDeferredBufferCap looks like SetBufferCap(),
// but sets a capacity of those encoded []byte entries, that is tried to be sent,
// while CI_WriterHttp is temporary disabled.
//
// Usually, you set this bigger than value used in SetBufferCap(),
// to store your entries when your provider is not available you don't know why,
// because exactly in this buffer all those entries that are tried to be processed,
// while your provider is not available, will be stored.
//
// If is set to 0, the entries will be discarded
// when CI_WriterHttp is temporary disabled.
//
// Read p.1, p.3, p.6 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// Allowed range: [0..8'388'608].
// Do not overwrite default value if you're not understand how much RAM
// your log entries consume per item, and what RAM consumption will be
// at the upper bound.
// Default: 16384.
func (dw *CI_WriterHttp) SetDeferredBufferCap(cap uint32) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		if cap <= uint32(1<<23) {
			dw.deferredEntriesBufferLen = &cap
		}
	})
}

// SetWorkersNum sets how much goroutines will be spawned
// to handle all passed encoded []byte entries and send them to your provider.
//
// Read p.1 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// Allowed range: [1..32].
// Set high values only if there is a really high throughput is required.
// Default: 2. Recommended: [1..4].
func (dw *CI_WriterHttp) SetWorkersNum(num uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		if num >= 1 && num <= 32 {
			dw.workerNum = num
		}
	})
}

// SetWorkerAutoFlushDelay sets how often accumulated []byte entries will be sent
// to your provider, even if their buffer is not full.
//
// Read p.3, p.5 of CI_WriterHttp doc for more info.
//
// Does nothing, if CI_WriterHttp already running, stopped or disabled
// (Write() has been called at least once).
//
// Allowed range: [100ms..24h].
// Default: 10s.
func (dw *CI_WriterHttp) SetWorkerAutoFlushDelay(delay time.Duration) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		if delay >= 100*time.Microsecond && delay <= 24*time.Hour {
			dw.workerFlushDelay = delay
		}
	})
}

// AddBefore sets the data that will be added to the encoded entries pack's buffer
// before the first encoded entry is added.
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "[" (JSON beginning of array char).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddBefore(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataBefore = data
	})
}

// AddAfter sets the data that will be added to the encoded entries pack's buffer
// after the last encoded entry is added.
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "]" (JSON ending of array char).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddAfter(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataAfter = data
	})
}

// AddBetween sets the data that will be added to the encoded entries pack's buffer
// between encoded log entries (but neither before first nor after last).
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "," (JSON separator of objects inside an array).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddBetween(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataBetween = data
	})
}

// AddBeforeS is the same as AddBefore() but accepts string instead of []byte,
// doing no-copy conversion.
func (dw *CI_WriterHttp) AddBeforeS(data string) *CI_WriterHttp {
	return dw.AddBefore(ekastr.S2B(data))
}

// AddAfterS is the same as AddAfter() but accepts string instead of []byte,
// doing no-copy conversion.
func (dw *CI_WriterHttp) AddAfterS(data string) *CI_WriterHttp {
	return dw.AddAfter(ekastr.S2B(data))
}

// AddBetweenS is the same as AddBetween() but accepts string instead of []byte,
// doing no-copy conversion.
func (dw *CI_WriterHttp) AddBetweenS(data string) *CI_WriterHttp {
	return dw.AddBetween(ekastr.S2B(data))
}

// AddBeforeAfterBetween replaces all of AddBefore(), AddAfter(), AddBetween() calls.
//
// It accepts 1 or 3 arguments.
// If 3 arguments is provided, they passed to corresponding calls.
// If 1 argument is passed, and it's length is multiple of 3, it will be split
// to the equal length 3 pieces, and they are passed to corresponding calls.
//
// All other variants of arguments are ignored and does no-op.
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddBeforeAfterBetween(args ...[]byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		switch l := len(args); {
		case l == 1 && len(args[0]) > 0 && len(args[0])%3 == 0:
			l = len(args[0]) / 3
			dw.dataBefore, dw.dataAfter, dw.dataBetween =
				args[0][:l], args[0][l:l*2], args[0][l*2:]
		case l == 3:
			dw.dataBefore, dw.dataAfter, dw.dataBetween =
				args[0], args[1], args[2]
		default:
			return
		}
	})
}

// AddBeforeAfterBetweenS is the same as AddBeforeAfterBetweenS()
// but accepts strings instead of [][]byte, doing no-copy conversion.
func (dw *CI_WriterHttp) AddBeforeAfterBetweenS(args ...string) *CI_WriterHttp {
	var args_ [][]byte
	if len(args) > 0 {
		args_ = make([][]byte, len(args))
		for i, arg := range args {
			args_[i] = ekastr.S2B(arg)
		}
	}
	return dw.AddBeforeAfterBetween(args_...)
}

// Ping checks whether provider settings are correct and connection can be established.
//
// Keep in mind, that if you do not do ping by yourself, the CI_WriterHttp will try
// to ping when Write() will be called first time (initialization).
// But if you call method Ping(), there will be no internal Ping() call
// at the initialization.
// You may prefer explicit Ping() call, because here you may specify some callbacks
// 'cb' that are applied to the HTTP request before it will be send.
//
// WARNING!
// If you call Ping() by yourself and it returns an error, you may change something
// (provider settings, etc), and then try to ping again.
// But if Ping() is called at the initialization, the CI_WriterHttp can not be recover.
func (dw *CI_WriterHttp) Ping(cb ...func(req *fasthttp.Request)) *ekaerr.Error {
	switch {

	case dw == nil:
		return ekaerr.IllegalState.
			New("CI_WriterHttp: writer is nil (not initialized)").
			Throw()
	}

	return dw.ping(true, cb)
}

// Write sends 'p' to the internal entries being processed buffer and returns
// len(p) and nil if 'p' has been successfully queued.
//
// Initializes CI_WriterHttp object if it's not. If initialization once failed,
// the CI_WriterHttp can not be used anymore.
//
// Returned errors:
// - nil: OK, 'p' has been queued.
// - ErrWriterIsNil: CI_WriterHttp receiver is nil.
// - ErrWriterDisabled: CI_WriterHttp is stopped and will never start again.
// - ErrWriterBufferFull: Internal CI_WriterHttp's buffer of processed entries
//   is full. Next time set bigger buffer's length using SetBufferCap().
func (dw *CI_WriterHttp) Write(p []byte) (n int, err error) {
	switch {

	case dw == nil:
		return -1, ErrWriterIsNil

	case len(p) == 0:
		return 0, nil

	case !dw.canWrite():
		return -1, ErrWriterDisabled
	}

	select {

	case dw.entries <- p:
		return len(p), nil

	default:
		atomic.AddUint64(&dw.entriesCompletelyLostCounter, 1)
		return -1, ErrWriterBufferFull
	}
}
