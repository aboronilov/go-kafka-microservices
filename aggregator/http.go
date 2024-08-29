package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err  error
}

// implements Error interface
func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"message": apiErr.Error()})
			}
		}
	}
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	return &HTTPMetricHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
			Name:      "aggregator",
		}),
		errCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
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

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error

		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took":   time.Since(start).Seconds(),
				"path":   r.URL.Path,
				"method": r.Method,
				"ip":     r.RemoteAddr,
				"err":    err,
			})
			h.reqCounter.Inc()
			h.reqLatency.Observe(time.Since(start).Seconds())
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		if err != nil {
			logrus.Error(err)
		}

		return err
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("method not allowed: %s", r.Method),
			}
		}

		invoiceID := r.URL.Query().Get("obu")
		if invoiceID == "" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invoice ID is required"),
			}
			// return writeJSON(
			// 	w,
			// 	http.StatusBadRequest,
			// 	map[string]string{"error": "invoice ID is required"},
			// )
		}

		id, err := strconv.Atoi(invoiceID)
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid invoice ID: %s", invoiceID),
			}
		}

		invoice, err := svc.CalculateInvoice(id)
		if err != nil {
			return writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			// writeJSON(
			// 	w,
			// 	http.StatusBadRequest,
			// 	map[string]string{"error": "Method not allowed:"})
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("method not allowed"),
			}

		}

		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			// writeJSON(
			// 	w,
			// 	http.StatusInternalServerError,
			// 	map[string]string{"error": err.Error()})
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("wrong distance format"),
			}
		}

		if err := svc.AggregateDistance(dist); err != nil {
			return writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})

		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dist)

		return nil
	}
}
