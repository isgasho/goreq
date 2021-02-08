package header

import (
	"net/http"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

func ContentType(value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderContentType, value)
			return next(response, request)
		}
	}
}

func Accept(value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderAccept, value)
			return next(response, request)
		}
	}
}

func UserAgent(value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderUserAgent, value)
			return next(response, request)
		}
	}
}

func Referer(value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderReferer, value)
			return next(response, request)
		}
	}
}

func Origin(value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderOrigin, value)
			return next(response, request)
		}
	}
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func Add(key, value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Add(key, value)
			return next(response, request)
		}
	}
}

// Set sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func Set(key, value string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(key, value)
			return next(response, request)
		}
	}
}

// SetMap sets a map of headers represented by key-value pair.
func SetHeader(header http.Header) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			for k, v := range header {
				for _, vv := range v {
					request.Header.Set(k, vv)
				}
			}
			return next(response, request)
		}
	}
}

// SetMap sets a map of headers represented by key-value pair.
func SetMap(headers map[string]string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			for k, v := range headers {
				request.Header.Set(k, v)
			}
			return next(response, request)
		}
	}
}

// Del deletes the header fields associated with key.
func Del(key string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Del(key)
			return next(response, request)
		}
	}
}

// Del deletes all headers
func DelAll() wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header = make(http.Header)
			return next(response, request)
		}
	}
}
