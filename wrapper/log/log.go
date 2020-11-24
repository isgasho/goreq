package log

import (
	"fmt"

	"github.com/aiscrm/goreq"
)

func Dump() goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp, opts goreq.CallOptions) error {
			next(req, resp, opts)
			fmt.Println(resp.Dump())
			return nil
		}
	}
}
