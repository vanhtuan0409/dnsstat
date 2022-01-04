package main

import (
	"log"
	"sync"
)

type statistic struct {
	sync.RWMutex
}

type query struct {
	domain string
	typ    string
}

func newStatistic() *statistic {
	ret := new(statistic)
	return ret
}

func (s *statistic) observe(q *query) {
	log.Printf("domain: %s - type: %s", q.domain, q.typ)
}

func (s *statistic) reset() error {
	return nil
}
