package main

import (
	"log"

	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
	"google.golang.org/protobuf/proto"
)

func worker(ch <-chan []byte, stat *statistic, conf *config) {
	for frame := range ch {
		queries, err := parseQueries(frame)
		if err != nil {
			log.Printf("[ERR] parse frame. ERR: %+v", err)
			continue
		}

		for _, q := range queries {
			if !isDomainBlackListed(q.domain, conf.IgnoreDomains) {
				stat.observe(q)
			}
		}
	}
}

func isDomainBlackListed(query string, blacklists []DomainMatcher) bool {
	ret := false
	for _, m := range blacklists {
		if m.IsMatch(query) {
			ret = true
			break
		}
	}
	return ret
}

func parseQueries(frame []byte) ([]*query, error) {
	ret := []*query{}
	var data dnstap.Dnstap
	if err := proto.Unmarshal(frame, &data); err != nil {
		return ret, err
	}
	if *data.Type != dnstap.Dnstap_MESSAGE {
		return ret, nil
	}
	if data.Message.QueryMessage == nil {
		return ret, nil
	}

	msg := new(dns.Msg)
	if err := msg.Unpack(data.Message.QueryMessage); err != nil {
		return ret, err
	}
	if msg.Question == nil {
		return ret, nil
	}

	for _, q := range msg.Question {
		ret = append(ret, &query{
			domain: q.Name,
			typ:    dns.Type(q.Qtype).String(),
		})
	}

	return ret, nil
}
