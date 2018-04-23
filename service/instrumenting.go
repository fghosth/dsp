package service

import (
	_ "fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"jvole.com/dsp/rabbitmq"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram

	next DSPBidder
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s DSPBidder) DSPBidder {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}
func (s *instrumentingService) SyncIndex() []byte {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Bidder").Add(1)
		s.requestLatency.With("method", "Bidder").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.SyncIndex()
}

func (s *instrumentingService) StopBidByCID(cid uint32) bool {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Bidder").Add(1)
		s.requestLatency.With("method", "Bidder").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.StopBidByCID(cid)
}

func (s *instrumentingService) Bidder(data []byte, host string) string {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Bidder").Add(1)
		s.requestLatency.With("method", "Bidder").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.Bidder(data, host)
}

func (s *instrumentingService) ADXNotify(notify rabbitmq.WinNotify) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ADXNotify").Add(1)
		s.requestLatency.With("method", "ADXNotify").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.ADXNotify(notify)
}
