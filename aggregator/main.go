package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.Info("Invoice service started")
	httpListenAddr := flag.String("httpAddr", ":4000", "the listen port of HTTP transport")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen port of GRPC transport")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	svc = NewMetricsMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCtransport(*grpcListenAddr, svc))
	}()

	log.Fatal(makeHTTPTransport(*httpListenAddr, svc))
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
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

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoiceID := r.URL.Query().Get("obu")
		if invoiceID == "" {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invoice ID is required"})
			return
		}

		id, err := strconv.Atoi(invoiceID)
		if err != nil {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid invoice ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(id)
		if err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
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
