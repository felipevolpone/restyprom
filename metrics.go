package restyprom

import (
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	responseTimePerURL    *prometheus.HistogramVec
	responseStatusCounter *prometheus.CounterVec
	successCounterByURL   *prometheus.CounterVec
	failuresCounterByURL  *prometheus.CounterVec
	once                  sync.Once
)

const (
	ns = "resty_prom"
)

func init() {
	responseTimePerURL = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "response_time_per_url",
			Help:      "Response time of called URLs",
		},
		[]string{"url"},
	)

	responseStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "total_status_per_url",
			Help:      "Total of calls per URL and status code",
		},
		[]string{"status", "url"},
	)

	successCounterByURL = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "total_success_per_url",
			Help:      "Total of success calls per URL",
		},
		[]string{"url"},
	)

	failuresCounterByURL = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "total_failures_per_url",
			Help:      "Total of failure calls per URL",
		},
		[]string{"url"},
	)
}

func NewBasicClient() *resty.Client {
	c := resty.New()
	c.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		Collect(r)
		return nil
	})
	return c
}

func registerMetrics(prom prometheus.Registerer) {
	once.Do(func() {
		prom.MustRegister(
			responseTimePerURL,
			responseStatusCounter,
			successCounterByURL,
			failuresCounterByURL,
		)
	})
}

func Init() {
	registerMetrics(prometheus.DefaultRegisterer)
}

func InitWithRegister(r prometheus.Registerer) {
	registerMetrics(r)
}

func Collect(r *resty.Response) {
	if r.IsSuccess() {
		successCounterByURL.With(prometheus.Labels{"url": r.Request.URL}).Inc()
	}

	if r.IsError() {
		failuresCounterByURL.With(prometheus.Labels{"url": r.Request.URL}).Inc()
	}

	responseTimePerURL.With(prometheus.Labels{"url": r.Request.URL}).Observe(r.Time().Seconds())

	responseStatusCounter.With(prometheus.Labels{
		"url":    r.Request.URL,
		"status": fmt.Sprintf("%d", r.RawResponse.StatusCode),
	}).Inc()
}
