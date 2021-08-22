package trace

import (
	"net/http"
	"net/http/httputil"

	"github.com/opentracing/opentracing-go/log"

	"github.com/aiscrm/goreq/wrapper"
	"github.com/opentracing/opentracing-go"
)

func Trace(operationNames ...string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			reqURI := request.RequestURI
			if reqURI == "" {
				reqURI = request.URL.RequestURI()
			}
			operationName := "goreq:" + request.URL.Host + ":" + reqURI
			if len(operationNames) > 0 {
				operationName = operationNames[0]
			}
			span, _ := opentracing.StartSpanFromContext(request.Context(), operationName)
			defer span.Finish()
			err := next(response, request)
			if err != nil {
				span.LogFields(log.Error(err))
				span.SetTag("error", true)
			}
			reqData, _ := httputil.DumpRequest(request, true)
			span.LogFields(log.String("request.body", string(reqData)))
			respData, _ := httputil.DumpResponse(response, true)
			span.LogFields(log.String("response.body", string(respData)))
			span.SetTag("component", "goreq")
			span.SetTag("span.kind", "client")
			span.SetTag("http.url", reqURI)
			span.SetTag("http.method", request.Method)
			span.SetTag("http.status_code", response.Status)
			span.SetTag("peer.hostname", request.URL.Host)
			return err
		}
	}
}
