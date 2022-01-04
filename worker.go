package main

import (
	"fmt"
	"log"
	"strings"

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
			isBlacklist := false
			for _, root := range conf.IgnoreRootDomains {
				if strings.HasSuffix(q.domain, fmt.Sprintf("%s.", root)) {
					isBlacklist = true
					break
				}
			}

			if !isBlacklist {
				stat.observe(q)
			}
		}
	}
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
