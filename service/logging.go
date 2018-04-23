package service

import (
	_ "fmt"
	"path"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/rabbitmq"
)

type loggingService struct {
	logger *log.Logger
	next   DSPBidder
}

func NewLoggingService(logger *log.Logger, s DSPBidder) DSPBidder {
	return &loggingService{logger, s}
}

func (s *loggingService) SyncIndex() []byte {
	defer func(begin time.Time) {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Info(*s.logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"took", time.Since(begin).Nanoseconds()/1000, //微妙
		)
	}(time.Now())
	return s.next.SyncIndex()
}

func (s *loggingService) StopBidByCID(cid uint32) bool {
	defer func(begin time.Time) {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Info(*s.logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"took", time.Since(begin).Nanoseconds()/1000, //微妙
		)
	}(time.Now())
	return s.next.StopBidByCID(cid)
}

func (s *loggingService) Bidder(data []byte, host string) string {
	defer func(begin time.Time) {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Info(*s.logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"took", time.Since(begin).Nanoseconds()/1000, //微妙
		)
	}(time.Now())
	return s.next.Bidder(data, host)
}

func (s *loggingService) ADXNotify(notify rabbitmq.WinNotify) error {
	defer func(begin time.Time) {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Info(*s.logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"took", time.Since(begin).Nanoseconds()/1000,
		)
	}(time.Now())

	return s.next.ADXNotify(notify)
}
