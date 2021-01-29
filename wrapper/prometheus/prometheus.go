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
	Gatherer   prometheus.Gatherer
	Objectives map[float64]float64
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		NameSpace:  "goreq",
		Registerer: prometheus.DefaultRegisterer,
		Gatherer:   prometheus.DefaultGatherer,
		Objectives: map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.98: 0.001, 1: 0},
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

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

func Gatherer(gatherer prometheus.Gatherer) Option {
	return func(options *Options) {
		options.Gatherer = gatherer
	}
}

func Objectives(objectives map[float64]float64) Option {
	return func(options *Options) {
		options.Objectives = objectives
	}
}

func Prometheus(opts ...Option) wrapper.CallWrapper {
	options := NewOptions(opts...)
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
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  options.NameSpace,
			Name:       "latency_microseconds",
			Objectives: options.Objectives,
			Help:       "Request latencies in microseconds, partitioned by host and uri",
		},
		[]string{
			"host",
			"uri",
		},
	)
	options.Registerer.MustRegister(counter)
	options.Registerer.MustRegister(summary)
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				us := v * 1000000 // make microseconds
				summary.WithLabelValues(request.RequestURI).Observe(us)
			}))
			err := next(response, request)
			counter.WithLabelValues(request.Host, request.RequestURI, strconv.Itoa(response.StatusCode)).Inc()
			timer.ObserveDuration()
			return err
		}
	}
	return nil
}
