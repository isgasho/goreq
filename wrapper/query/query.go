package query

import (
	"github.com/aiscrm/goreq"
	"github.com/aiscrm/goreq/util"
)

// Set sets the query param key and value.
// It replaces any existing values.
func Set(key string, value interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			query := req.Request.URL.Query()
			query.Set(key, util.ToString(value))
			req.Request.URL.RawQuery = query.Encode()
			return next(req, resp, opts)
		}
	}
}

// Add adds the query param value to key.
// It appends to any existing values associated with key.
func Add(key string, value interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			query := req.Request.URL.Query()
			query.Add(key, util.ToString(value))
			req.Request.URL.RawQuery = query.Encode()
			return next(req, resp, opts)
		}
	}
}

// Del deletes the query param values associated with key.
func Del(key string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			query := req.Request.URL.Query()
			query.Del(key)
			req.Request.URL.RawQuery = query.Encode()
			return next(req, resp, opts)
		}
	}
}

// DelAll deletes all the query params.
func DelAll() goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.URL.RawQuery = ""
			return next(req, resp, opts)
		}
	}
}

// SetMap sets a map of query params by key-value pair.
func SetMap(params map[string]interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			query := req.Request.URL.Query()
			for k, v := range params {
				query.Set(k, util.ToString(v))
			}
			req.Request.URL.RawQuery = query.Encode()
			return next(req, resp, opts)
		}
	}
}

// AddMap add a map of query params by key-value pair.
func AddMap(params map[string]interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			query := req.Request.URL.Query()
			for k, v := range params {
				query.Add(k, util.ToString(v))
			}
			req.Request.URL.RawQuery = query.Encode()
			return next(req, resp, opts)
		}
	}
}

//func Struct(v interface{}) goreq.CallWrapper {
//	return func(next goreq.CallFunc) goreq.CallFunc {
//		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
//			return next(req, resp, opts)
//		}
//	}
//}
