package form

import (
	"bytes"
	"mime/multipart"
	"strings"

	"github.com/aiscrm/goreq"
	"github.com/aiscrm/goreq/util"
)

// Set sets the query param key and value.
// It replaces any existing values.
func Set(key string, value interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.PostForm.Set(key, util.ToString(value))
			generate(req)
			return next(req, resp, opts)
		}
	}
}

// Add adds the query param key and value.
// It replaces any existing values.
func Add(key string, value interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			req.Request.PostForm.Add(key, util.ToString(value))
			generate(req)
			return next(req, resp, opts)
		}
	}
}

// SetMap sets a map of query params by key-value pair.
func SetMap(params map[string]interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			for key, value := range params {
				req.Request.PostForm.Set(key, util.ToString(value))
				generate(req)
			}
			return next(req, resp, opts)
		}
	}
}

// AddMap adds a map of query params by key-value pair.
func AddMap(params map[string]interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			for key, value := range params {
				req.Request.PostForm.Add(key, util.ToString(value))
				generate(req)
			}
			return next(req, resp, opts)
		}
	}
}

func generate(req *goreq.Req) error {
	if strings.Contains(req.Request.Header.Get(util.HeaderContentType), util.HeaderContentTypeMultipart) {
		data := new(bytes.Buffer)
		bodyWriter := multipart.NewWriter(data)

		if err := util.ToMultipart(req.Request, bodyWriter); err != nil {
			return err
		}

		_ = bodyWriter.Close()
		return util.SetBinary(req.Request, bytes.NewReader(data.Bytes()))
	} else {
		req.Request.Header.Set(util.HeaderContentType, util.HeaderContentTypeForm)
		return util.SetBinary(req.Request, bytes.NewReader([]byte(req.Request.PostForm.Encode())))
	}
}
