package goreq

import (
	"net/http"
	"net/url"
)

func Get(rawURL string) *Req {
	req := New()
	req.WithMethod(http.MethodGet)
	u, err := url.Parse(rawURL)
	if err != nil {
		req.Error = err
	}
	req.Request.URL = u
	return req
}

// Post return a post request
func Post(rawURL string) *Req {
	req := New()
	req.WithMethod(http.MethodPost)
	u, err := url.Parse(rawURL)
	if err != nil {
		req.Error = err
	}
	req.Request.URL = u
	return req
}
