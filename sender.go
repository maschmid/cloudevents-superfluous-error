package main

import (
	"fmt"
	"context"
	"net/http"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"log"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	const count = 64
	// 30 kB
	const size = 30 * 1024
	var ids uint64

	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("error creating default client: %v", err)
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), "http://127.0.0.1:8080")

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := cloudevents.NewEvent()

			id := atomic.AddUint64(&ids, 1)
			e.SetID(fmt.Sprintf("%d", id))
			e.SetSource("sender")
			e.SetType("myevent")

			// 30kB
			data := make([]byte, size, size)
			for c := 0; c < size; c++ {
				// anything...
				data[c] = 42
			}

			e.SetData("", data)

			result := c.Send(ctx, e)
			if !cloudevents.IsACK(result) {
				s := fmt.Sprintf("error sending event %d/%d: %v\n", i, count, result)
				log.Print(s)
			}
		}()
	}

	wg.Wait()

	resp, err := http.Get("http://127.0.0.1:8080/report")
	if err != nil {
		log.Fatalf("error getting report")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading report")
	}

	log.Print(string(body))
}
