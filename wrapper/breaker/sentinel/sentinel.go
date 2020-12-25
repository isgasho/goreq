package sentinel

import (
	"net/http"

	"github.com/aiscrm/goreq"
	sentinel "github.com/alibaba/sentinel-golang/api"
)

func Breaker(entryOpts ...sentinel.EntryOption) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			e, b := sentinel.Entry(defaultKey(req.Request), entryOpts...)
			if b != nil {
				// 请求被流控，可以从 BlockError 中获取限流详情
				// block 后不需要进行 Exit()
				return b
			} else {
				// 务必保证业务逻辑结束后 Exit
				defer e.Exit()
				// 请求可以通过，在此处编写您的业务逻辑
				next(req, resp, opts)
			}
			return nil
		}
	}
}

var defaultKey = func(request *http.Request) string {
	return request.RequestURI
}
