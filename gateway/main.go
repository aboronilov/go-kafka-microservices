package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/aboronilov/go-kafka-microservices/aggregator/client"
	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "http gateway port")
	flag.Parse()

	var (
		client     = client.NewClient("http://localhost:3000")
		invHandler = newInvoiceHandler(*client)
	)

	http.HandleFunc("/invoice", makeAPIfunc(invHandler.handleGetInvoice))
	logrus.Infof("Gateway is running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.HTTPClient
}

func newInvoiceHandler(client client.HTTPClient) *InvoiceHandler {
	return &InvoiceHandler{client: client}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	invoice, err := h.client.GetInvoice(context.Background(), 1)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, map[string]*types.Invoice{"invoice": invoice})
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIfunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
