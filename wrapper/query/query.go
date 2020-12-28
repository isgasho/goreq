package query

import (
	"net/http"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

// Set sets the query param key and value.
// It replaces any existing values.
func Set(key string, value interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			query := request.URL.Query()
			query.Set(key, util.ToString(value))
			request.URL.RawQuery = query.Encode()
			return next(response, request)
		}
	}
}

// Add adds the query param value to key.
// It appends to any existing values associated with key.
func Add(key string, value interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			query := request.URL.Query()
			query.Add(key, util.ToString(value))
			request.URL.RawQuery = query.Encode()
			return next(response, request)
		}
	}
}

// Del deletes the query param values associated with key.
func Del(key string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			query := request.URL.Query()
			query.Del(key)
			request.URL.RawQuery = query.Encode()
			return next(response, request)
		}
	}
}

// DelAll deletes all the query params.
func DelAll() wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.URL.RawQuery = ""
			return next(response, request)
		}
	}
}

// SetMap sets a map of query params by key-value pair.
func SetMap(params map[string]interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			query := request.URL.Query()
			for k, v := range params {
				query.Set(k, util.ToString(v))
			}
			request.URL.RawQuery = query.Encode()
			return next(response, request)
		}
	}
}

// AddMap add a map of query params by key-value pair.
func AddMap(params map[string]interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			query := request.URL.Query()
			for k, v := range params {
				query.Add(k, util.ToString(v))
			}
			request.URL.RawQuery = query.Encode()
			return next(response, request)
		}
	}
}

//func Struct(v interface{}) goreq.CallWrapper {
//	return func(next goreq.CallFunc) goreq.CallFunc {
//		return func(response *http.Response, request *http.Request) error {
//			return next(response, request)
//		}
//	}
//}
