package url

import (
	"net/http"
	"net/url"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

// URL set request raw url
// 建议在client使用URL方法，在Req使用AddPath/Path方法。好处：
// - 不同运行环境时api的域名不同时，方便根据环境进行切换调用的接口地址
// - 当api地址发生变化时可以容易的一次性更改api的地址
func URL(rawURL string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			u, err := url.Parse(rawURL)
			if err != nil {
				return err
			}
			if request.URL == nil {
				request.URL = u
			} else {
				request.URL.Scheme = u.Scheme
				request.URL.Host = u.Host

				if u.Path != "" && u.Path != "/" {
					request.URL.Path = u.Path + request.URL.Path
				}
			}
			return next(response, request)
		}
	}
}

// Path set request raw url
func Path(path string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			if request.URL == nil {
				request.URL = new(url.URL)
			}
			request.URL.Path = util.NormalizePath(path)
			return next(response, request)
		}
	}
}

// AddPath add Request path
func AddPath(path string) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			if request.URL == nil {
				request.URL = new(url.URL)
			}
			request.URL.Path += util.NormalizePath(path)
			return next(response, request)
		}
	}
}
