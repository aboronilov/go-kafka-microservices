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
	svc = NewLogMiddleware(svc)

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
			writwJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Methot not allowed:"})
			return
		}

		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			writwJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(dist); err != nil {
			writwJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dist)
	}
}

func writwJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
