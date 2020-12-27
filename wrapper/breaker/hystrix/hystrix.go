package hystrix

import (
	"net/http"

	"github.com/aiscrm/goreq"
	"github.com/aiscrm/goreq/wrapper/breaker"

	"github.com/afex/hystrix-go/hystrix"
)

func Breaker(opts ...Option) goreq.CallWrapper {
	options := Options{
		KeyFunc: func(request *http.Request) string {
			return request.RequestURI
		},
		Commands: make(map[string]hystrix.CommandConfig),
	}
	for _, opt := range opts {
		opt(&options)
	}

	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			var doError error
			hystrix.Do(options.KeyFunc(req.Request), func() error {
				return next(req, resp)
			}, func(err error) error {
				doError = breaker.CircuitError{
					Message: err.Error(),
				}
				return nil
			})
			return doError
		}
	}
}
