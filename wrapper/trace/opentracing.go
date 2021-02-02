package trace

import (
	"net/http"

	opentracinglog "github.com/opentracing/opentracing-go/log"

	"github.com/aiscrm/goreq/wrapper"
	"github.com/opentracing/opentracing-go"
)

func Trace() wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			span, _ := opentracing.StartSpanFromContext(request.Context(), "goreq")
			defer span.Finish()
			err := next(response, request)
			if err != nil {
				span.LogFields(opentracinglog.Error(err))
				span.SetTag("error", true)
			}
			span.LogFields(opentracinglog.String("req.host", request.URL.Host))
			span.LogFields(opentracinglog.String("req.url", request.URL.Path))
			span.LogFields(opentracinglog.String("resp.status", response.Status))
			return err
		}
	}
}
