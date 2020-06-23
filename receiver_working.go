package main

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	event "github.com/cloudevents/sdk-go/v2/event"
)

var count uint64

func receive(e event.Event) error {
	log.Printf("received %s", e.ID())
	atomic.AddUint64(&count, 1)

	return nil
}

func report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("count = %d\n", atomic.LoadUint64(&count))))
}

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	router := http.NewServeMux()
	p.Handler = router

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	router.HandleFunc("/report", report)

	log.Printf("will listen on :8080\n")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, receive))
}
