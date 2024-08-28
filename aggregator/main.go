package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.Info("Invoice service started")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_PORT")
		httpListenAddr = os.Getenv("AGG_HTTP_PORT")
	)

	svc = NewLogMiddleware(svc)
	svc = NewMetricsMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCtransport(grpcListenAddr, svc))
	}()

	log.Fatal(makeHTTPTransport(httpListenAddr, svc))
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port", listenAddr)

	aggregateMetric := NewHTTPMetricHandler("aggregate")
	invoiceMetric := NewHTTPMetricHandler("invoice")

	http.HandleFunc("/aggregate", aggregateMetric.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invoiceMetric.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAddr, nil)
}

func makeGRPCtransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport is running on port", listenAddr)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer func() {
		fmt.Println("GRPC transport is stopping")
		ln.Close()
	}()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGrpcAggregatorServer(svc))
	return server.Serve(ln)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Methot not allowed:"})
			return
		}

		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(dist); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dist)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")

	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("Unknown store type: %s", storeType)
		return nil
	}
}
