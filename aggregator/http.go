package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	return &HTTPMetricHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
			Name:      "aggregator",
		}),
		reqLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
			Name:      "aggregator",
			Help:      "Request latency in milliseconds",
			Buckets:   []float64{10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000},
		}),
	}
}

func (h *HTTPMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start).Seconds(),
				"path":   r.URL.Path,
				"method": r.Method,
				"ip":     r.RemoteAddr,
			})
			h.reqLatency.Observe(time.Since(start).Seconds())
		}(time.Now())
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Methot not allowed:"})
			return
		}

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
