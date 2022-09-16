package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. Розробити http сервер з двома ендпоінтами
// / - відповідає 200 з вірогідністю 33%, у всіх інших випадках 503
// /metrics - віддає prometheus метрики по кількості запитив з відповіддю 200 і 503

var (
	OKMetricProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ok200_processed_ops_total",
		Help: "The total number of processed 200 status code",
	})
	SAMetricProcess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sa503_processed_ops_total",
		Help: "The total number of processed 503 status code",
	})
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Starting api server")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		n := rand.Intn(100)
		if n < 33 {
			w.WriteHeader(200)
			OKMetricProcess.Inc()
			fmt.Fprintf(w, "OK")
			log.Println("responce 200 OK")
		} else {
			w.WriteHeader(503)
			SAMetricProcess.Inc()
			fmt.Fprintf(w, "Service Unavailable")
			log.Println("responce 503 Service Unavailable")
		}
	default:
		fmt.Fprintf(w, "Only GET supported.")
	}
}
