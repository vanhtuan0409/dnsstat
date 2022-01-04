package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"

	dnstap "github.com/dnstap/golang-dnstap"
)

type config struct {
	DataFile string
	SockFile string
	Port     int
	Worker   int
	Bufsize  int
	Topk     int
	HttpPort int

	IgnoreRootDomainsRaw string
	IgnoreRootDomains    []string

	Listener   net.Listener
	Input      dnstap.Input
	HttpServer *http.Server
}

func parseConfig() *config {
	ret := new(config)
	flag.StringVar(&ret.DataFile, "data", "", "Path to dnstap data file (used for testing purpose)")
	flag.IntVar(&ret.Port, "port", 0, "TCP port to receive dnstap data")
	flag.IntVar(&ret.Worker, "worker", 0, "Number of worker")
	flag.IntVar(&ret.Bufsize, "buf", 100, "Channel buffer for receiving dnstap data")
	flag.StringVar(&ret.IgnoreRootDomainsRaw, "ignore-root", "", "Ignore root domain query (comma sep)")
	flag.IntVar(&ret.Topk, "topk", 10, "Number of top frequent domain to keep track")
	flag.IntVar(&ret.HttpPort, "http", 6385, "Port to bind http")
	flag.Parse()

	if err := ret.Validate(); err != nil {
		panic(err)
	}

	return ret
}

func (c *config) Validate() (err error) {
	if c.DataFile != "" {
		c.Input, err = dnstap.NewFrameStreamInputFromFilename(c.DataFile)
		if err != nil {
			return
		}
	}

	if c.SockFile != "" {
		c.Input, err = dnstap.NewFrameStreamSockInputFromPath(c.SockFile)
		if err != nil {
			return
		}
	}

	if c.Port != 0 {
		addr := fmt.Sprintf(":%d", c.Port)
		c.Listener, err = net.Listen("tcp", addr)
		if err != nil {
			return
		}
		c.Input = dnstap.NewFrameStreamSockInput(c.Listener)
	}

	if c.Input == nil {
		err = errors.New("No input provided")
		return
	}

	if c.Worker == 0 {
		c.Worker = runtime.NumCPU()
	}

	if c.IgnoreRootDomainsRaw != "" {
		c.IgnoreRootDomains = strings.Split(c.IgnoreRootDomainsRaw, ",")
	}

	if c.HttpPort != 0 {
		c.HttpServer = &http.Server{
			Addr: fmt.Sprintf(":%d", c.HttpPort),
		}
	}

	return
}

func (c *config) Close() error {
	if c.Listener != nil {
		c.Listener.Close()
	}
	return nil
}
