package main

import (
	"sync"

	"github.com/vanhtuan0409/dnsstat/internal/topk"
)

type statistic struct {
	stream *topk.Stream
	sync.RWMutex
}

type query struct {
	domain string
	typ    string
}

func newStatistic(conf *config) *statistic {
	ret := new(statistic)
	ret.stream = topk.New(conf.Topk)
	return ret
}

func (s *statistic) observe(q *query) {
	s.Lock()
	defer s.Unlock()
	s.stream.Insert(q.domain, 1)
}

func (s *statistic) reset() error {
	return nil
}
