package header

import (
	"net/http"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq"
)

func ContentType(value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderContentType, value)
			return next(req, resp)
		}
	}
}

func Accept(value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderAccept, value)
			return next(req, resp)
		}
	}
}

func UserAgent(value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderUserAgent, value)
			return next(req, resp)
		}
	}
}

func Referer(value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderReferer, value)
			return next(req, resp)
		}
	}
}

func Origin(value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderOrigin, value)
			return next(req, resp)
		}
	}
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func Add(key, value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Add(key, value)
			return next(req, resp)
		}
	}
}

// Set sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func Set(key, value string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(key, value)
			return next(req, resp)
		}
	}
}

// SetMap sets a map of headers represented by key-value pair.
func SetMap(headers map[string]string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			for k, v := range headers {
				req.Request.Header.Set(k, v)
			}
			return next(req, resp)
		}
	}
}

// Del deletes the header fields associated with key.
func Del(key string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Del(key)
			return next(req, resp)
		}
	}
}

// Del deletes all headers
func DelAll() goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header = make(http.Header)
			return next(req, resp)
		}
	}
}
