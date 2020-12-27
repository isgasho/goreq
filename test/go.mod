module github.com/aiscrm/goreq/test

go 1.15

replace (
	github.com/aiscrm/goreq => ../
	github.com/aiscrm/goreq/vo => ../vo
	github.com/aiscrm/goreq/wrapper/breaker/hystrix => ../wrapper/breaker/hystrix
)

require (
	github.com/aiscrm/goreq v0.0.0-00010101000000-000000000000
	github.com/aiscrm/goreq/vo v0.0.0-00010101000000-000000000000
	github.com/aiscrm/goreq/wrapper/breaker/hystrix v0.0.0-00010101000000-000000000000
)
