package hystrix

import (
	"net/http"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/wrapper/breaker"

	"github.com/afex/hystrix-go/hystrix"
)

func Breaker(opts ...Option) wrapper.CallWrapper {
	options := Options{
		KeyFunc: func(request *http.Request) string {
			return request.RequestURI
		},
		Commands: make(map[string]hystrix.CommandConfig),
	}
	for _, opt := range opts {
		opt(&options)
	}

	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			var circuitError error
			hystrix.Do(options.KeyFunc(request), func() error {
				return next(response, request)
			}, func(err error) error {
				circuitError = breaker.CircuitError{
					Message: err.Error(),
				}
				return nil
			})
			return circuitError
		}
	}
}
