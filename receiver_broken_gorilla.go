package main

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	event "github.com/cloudevents/sdk-go/v2/event"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, receive)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	router.Handle("/", h)
	router.HandleFunc("/report", report)

	log.Printf("will listen on :8080\n")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("unable to start http server, %s", err)
	}
}
