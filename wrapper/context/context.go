package context

import (
	"context"

	"github.com/aiscrm/goreq"
)

func Context(ctx context.Context) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.WithContext(ctx)
			return next(req, resp, opts)
		}
	}
}

func Value(key, val interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			ctx := context.WithValue(req.Request.Context(), key, val)
			req.Request.WithContext(ctx)
			return next(req, resp, opts)
		}
	}
}
