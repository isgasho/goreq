package goreq

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aiscrm/goreq/wrapper"
)

// Req represents a http request
type Req struct {
	Request *http.Request
	Error   error
	client  Client
	//body     []byte
	wrappers []wrapper.CallWrapper
}

// New return a empty request
func New() *Req {
	request := &http.Request{
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Method:     http.MethodGet,
		PostForm:   make(url.Values),
	}
	return &Req{
		client:  DefaultClient,
		Request: request,
		Error:   nil,
	}
}

func (r *Req) Clone() *Req {
	clone := new(Req)
	clone.client = r.client
	clone.Request = r.Request.Clone(r.Request.Context())
	clone.Error = r.Error
	return clone
}

func (r *Req) Use(wrappers ...wrapper.CallWrapper) *Req {
	r.wrappers = append(r.wrappers, wrappers...)
	return r
}

//WithClient with client
func (r *Req) WithClient(c Client) *Req {
	r.client = c
	//if c.Options().Endpoint != "" {
	//	r.Use(WithURL(c.Options().Endpoint))
	//}
	return r
}

// WithURL set request raw url
func (r *Req) WithURL(rawURL string) *Req {
	r.Request.URL, r.Error = url.Parse(rawURL)
	return r
}

func (r *Req) WithMethod(method string) *Req {
	r.Request.Method = method
	return r
}

// GetBody return request body
func (r *Req) GetBody() []byte {
	if r.Request.Body == nil {
		return []byte{}
	}
	body, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return []byte{}
	}
	r.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	return body
}

// GetURL return request context
func (r *Req) GetContext() context.Context {
	return r.Request.Context()
}

// GetURL return client
func (r *Req) GetClient() Client {
	return r.client
}

// Do is to call the request
func (r *Req) Do() *Resp {
	if r.client == nil {
		r.client = DefaultClient
	}
	return r.client.Do(r)
}
