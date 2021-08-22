module github.com/aiscrm/goreq/wrapper/trace

go 1.15

replace github.com/aiscrm/goreq => ../../

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/aiscrm/goreq v0.1.12
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
)
