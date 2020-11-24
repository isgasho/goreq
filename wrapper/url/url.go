package url

import (
	"net/url"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq"
)

// URL set request raw url
func URL(rawURL string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			u, err := url.Parse(rawURL)
			if err != nil {
				return err
			}
			if req.Request.URL == nil {
				req.Request.URL = u //new(url.URL)
			} else {
				req.Request.URL.Scheme = u.Scheme
				req.Request.URL.Host = u.Host

				if u.Path != "" && u.Path != "/" {
					req.Request.URL.Path = u.Path + req.Request.URL.Path
				}
			}
			return next(req, resp, opts)
		}
	}
}

// Path set request raw url
func Path(path string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			if req.Request.URL == nil {
				req.Request.URL = new(url.URL)
			}
			req.Request.URL.Path = util.NormalizePath(path)
			return next(req, resp, opts)
		}
	}
}

// AddPath add Request path
func AddPath(path string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			if req.Request.URL == nil {
				req.Request.URL = new(url.URL)
			}
			req.Request.URL.Path += util.NormalizePath(path)
			return next(req, resp, opts)
		}
	}
}
