package main

import (
	"time"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

type MetricsMiddleware struct {
	next           Aggregator
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0, 0.25, 0.5, 0.75, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0, 0.25, 0.5, 0.75, 1},
	})

	return &MetricsMiddleware{
		next:           next,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqCounterAgg.Inc()
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(dist)

	return err
}

func (m *MetricsMiddleware) CalculateInvoice(obu int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqCounterCalc.Inc()
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())

	invoice, err = m.next.CalculateInvoice(obu)

	return invoice, err
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("aggregating the distance")
	}(time.Now())

	return m.next.AggregateDistance(dist)
}

func (m *LogMiddleware) CalculateInvoice(obu int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if invoice != nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obu":      obu,
			"distance": distance,
			"amount":   amount,
		}).Info("calculating the invoice")
	}(time.Now())

	return m.next.CalculateInvoice(obu)
}
