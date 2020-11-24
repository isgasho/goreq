package cookie

import (
	"net/http"

	"github.com/aiscrm/goreq"
)

// Add adds a cookie to the request. Per RFC 6265 section 5.4, AddCookie does not
// attach more than one Cookie header field.
// That means all cookies, if any, are written into the same line, separated by semicolon.
func Add(cookie *http.Cookie) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.AddCookie(cookie)
			return next(req, resp, opts)
		}
	}
}

// DelAll deletes all the cookies by deleting the Cookie header field.
func DelAll() goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.Header.Del("Cookie")
			return next(req, resp, opts)
		}
	}
}
