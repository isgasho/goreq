package cookie

import (
	"net/http"

	"github.com/aiscrm/goreq/wrapper"
)

// Add adds a cookie to the request. Per RFC 6265 section 5.4, AddCookie does not
// attach more than one Cookie header field.
// That means all cookies, if any, are written into the same line, separated by semicolon.
func Add(cookie *http.Cookie) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.AddCookie(cookie)
			return next(response, request)
		}
	}
}

// AddCookies add multi cookies.
func AddCookies(cookies []*http.Cookie) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}
			return next(response, request)
		}
	}
}

// DelAll deletes all the cookies by deleting the Cookie header field.
func DelAll() wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Del("Cookie")
			return next(response, request)
		}
	}
}
