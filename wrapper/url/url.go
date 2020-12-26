package url

import (
	"net/url"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq"
)

// URL set request raw url
// 建议在client使用URL方法，在Req使用AddPath/Path方法。好处：
// - 不同运行环境时api的域名不同时，方便根据环境进行切换调用的接口地址
// - 当api地址发生变化时可以容易的一次性更改api的地址
func URL(rawURL string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			u, err := url.Parse(rawURL)
			if err != nil {
				return err
			}
			if req.Request.URL == nil {
				req.Request.URL = u
			} else {
				req.Request.URL.Scheme = u.Scheme
				req.Request.URL.Host = u.Host

				if u.Path != "" && u.Path != "/" {
					req.Request.URL.Path = u.Path + req.Request.URL.Path
				}
			}
			return next(req, resp)
		}
	}
}

// Path set request raw url
func Path(path string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			if req.Request.URL == nil {
				req.Request.URL = new(url.URL)
			}
			req.Request.URL.Path = util.NormalizePath(path)
			return next(req, resp)
		}
	}
}

// AddPath add Request path
func AddPath(path string) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			if req.Request.URL == nil {
				req.Request.URL = new(url.URL)
			}
			req.Request.URL.Path += util.NormalizePath(path)
			return next(req, resp)
		}
	}
}
