package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf := parseConfig()
	stat := newStatistic()
	ch := make(chan []byte, conf.Bufsize)
	go conf.Input.ReadInto(ch)
	log.Printf("Started server with config: %+v", conf)

	var wg sync.WaitGroup
	wg.Add(conf.Worker)
	for i := 0; i < conf.Worker; i++ {
		go func() {
			worker(ch, stat, conf)
			wg.Done()
		}()
	}

	<-ctx.Done()
	log.Println("Shutting down")
	conf.Close()
	close(ch)
	wg.Wait()
}
