module github.com/aiscrm/goreq/wrapper/breaker/hystrix

go 1.15

replace github.com/aiscrm/goreq => ../../../

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aiscrm/goreq v0.1.11
	github.com/smartystreets/goconvey v1.6.4 // indirect
)
