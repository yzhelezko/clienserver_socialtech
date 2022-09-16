package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// 2. Розробити http клієнт який працює в 2 потока, кожний потік робить запит один раз в секунду на http сервер,
// у випадку 503 робить 4 ретрая використовуючи backoff алгоритм (wait +5ms, +10ms, +20ms, +40ms) після кожної спроби.
// Response code після кожного запиту пишемо в лог

var backoffSchedule = []time.Duration{
	5 * time.Second,
	10 * time.Second,
	20 * time.Second,
	40 * time.Second,
}
var wg sync.WaitGroup

func main() {

	go runClient("first")
	go runClient("second")
	wg.Add(2)
	wg.Wait()
}

func runClient(run string) {

	for {
		var statusCode int
		for _, backoff := range backoffSchedule {
			statusCode = getIndex()
			log.Println(run, statusCode)
			if statusCode == 200 {
				break
			}

			fmt.Fprintf(os.Stderr, "Request error: %+v\n", 503)
			fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
			time.Sleep(backoff)
		}

		if statusCode != 200 {
			log.Println(run, "FATAL")
			//break
		}

		time.Sleep(time.Second)
	}
}

func getIndex() int {
	resp, err := http.Get("http://localhost:8000/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
