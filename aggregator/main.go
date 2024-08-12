package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Invoice service started")
	listenAddr := flag.String("listenaddr", ":3000", "the listen port of HTTP transport")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := svc.AggregateDistance(dist); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dist)
	}
}
