package goreq

import (
	"net/http"

	"github.com/aiscrm/goreq/wrapper/url"
)

func Get(rawURL string) *Req {
	return New().WithMethod(http.MethodGet).Use(url.URL(rawURL))
}

// Post return a post request
func Post(rawURL string) *Req {
	return New().WithMethod(http.MethodPost).Use(url.URL(rawURL))
}
