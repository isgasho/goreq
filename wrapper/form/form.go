package form

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

// Set sets the query param key and value.
// It replaces any existing values.
func Set(key string, value interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.PostForm.Set(key, util.ToString(value))
			generate(request)
			return next(response, request)
		}
	}
}

// Add adds the query param key and value.
// It replaces any existing values.
func Add(key string, value interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.PostForm.Add(key, util.ToString(value))
			generate(request)
			return next(response, request)
		}
	}
}

// SetMap sets a map of query params by key-value pair.
func SetMap(params map[string]interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			for key, value := range params {
				request.PostForm.Set(key, util.ToString(value))
				generate(request)
			}
			return next(response, request)
		}
	}
}

// AddMap adds a map of query params by key-value pair.
func AddMap(params map[string]interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			for key, value := range params {
				request.PostForm.Add(key, util.ToString(value))
				generate(request)
			}
			return next(response, request)
		}
	}
}

func generate(request *http.Request) error {
	if strings.Contains(request.Header.Get(util.HeaderContentType), util.HeaderContentTypeMultipart) {
		data := new(bytes.Buffer)
		bodyWriter := multipart.NewWriter(data)

		if err := util.ToMultipart(request, bodyWriter); err != nil {
			return err
		}

		_ = bodyWriter.Close()
		return util.SetBinary(request, bytes.NewReader(data.Bytes()))
	} else {
		request.Header.Set(util.HeaderContentType, util.HeaderContentTypeForm)
		return util.SetBinary(request, bytes.NewReader([]byte(request.PostForm.Encode())))
	}
}
