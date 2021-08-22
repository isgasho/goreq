module github.com/aiscrm/goreq/test

go 1.15

replace (
	github.com/aiscrm/goreq => ../
	github.com/aiscrm/goreq/vo => ../vo
	github.com/aiscrm/goreq/wrapper/breaker/hystrix => ../wrapper/breaker/hystrix
	github.com/aiscrm/goreq/wrapper/trace => ../wrapper/trace
)

require (
	github.com/aiscrm/goreq v0.1.12
	github.com/aiscrm/goreq/vo v0.0.0-00010101000000-000000000000
	github.com/aiscrm/goreq/wrapper/breaker/hystrix v0.0.0-00010101000000-000000000000
	github.com/aiscrm/goreq/wrapper/trace v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1 // indirect
)
