package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf := parseConfig()
	stat := newStatistic(conf)
	ch := make(chan []byte, conf.Bufsize)
	go conf.Input.ReadInto(ch)
	log.Printf("Started dnstap server with config: %+v", conf)

	var wg sync.WaitGroup
	wg.Add(conf.Worker)
	for i := 0; i < conf.Worker; i++ {
		go func() {
			worker(ch, stat, conf)
			wg.Done()
		}()
	}

	if conf.HttpServer != nil {
		http.HandleFunc("/", statHandler(stat))
		go conf.HttpServer.ListenAndServe()
		log.Printf("Started http server at :%d", conf.HttpPort)
	}

	<-ctx.Done()
	log.Println("Shutting down")
	conf.Close()
	close(ch)
	wg.Wait()
	if conf.HttpServer != nil {
		gracefulCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conf.HttpServer.Shutdown(gracefulCtx)
	}
}
