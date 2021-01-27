package prometheus

import (
	"net/http"
	"strconv"

	"github.com/aiscrm/goreq/wrapper"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	NameSpace  string
	Registerer prometheus.Registerer
}

type Option func(*Options)

func NameSpace(nameSpace string) Option {
	return func(options *Options) {
		options.NameSpace = nameSpace
	}
}

func Registerer(registerer prometheus.Registerer) Option {
	return func(options *Options) {
		options.Registerer = registerer
	}
}

func Prometheus(opts ...Option) wrapper.CallWrapper {
	options := Options{
		NameSpace:  "goreq",
		Registerer: prometheus.DefaultRegisterer,
	}
	for _, opt := range opts {
		opt(&options)
	}
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: options.NameSpace,
			Name:      "request_total",
			Help:      "Requests processed, partitioned by host, uri and status",
		},
		[]string{
			"host",
			"uri",
			"status",
		},
	)
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: options.NameSpace,
			Name:      "request_duration_seconds",
			Help:      "Request time in seconds, partitioned by host and uri",
		},
		[]string{
			"host",
			"uri",
		},
	)
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: options.NameSpace,
			Name:      "latency_microseconds",
			Help:      "Request latencies in microseconds, partitioned by host and uri",
		},
		[]string{
			"host",
			"uri",
		},
	)
	options.Registerer.Register(counter)
	options.Registerer.Register(histogram)
	options.Registerer.Register(summary)
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				us := v * 1000000 // make microseconds
				summary.WithLabelValues(request.RequestURI).Observe(us)
				histogram.WithLabelValues(request.Host, request.RequestURI).Observe(v)
			}))
			err := next(response, request)
			counter.WithLabelValues(request.Host, request.RequestURI, strconv.Itoa(response.StatusCode)).Inc()
			timer.ObserveDuration()
			return err
		}
	}
	return nil
}
