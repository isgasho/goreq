package context

import (
	"context"
	"net/http"

	"github.com/aiscrm/goreq/wrapper"
)

func Context(ctx context.Context) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.WithContext(ctx)
			return next(response, request)
		}
	}
}

func Value(key, val interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			ctx := context.WithValue(request.Context(), key, val)
			request.WithContext(ctx)
			return next(response, request)
		}
	}
}
