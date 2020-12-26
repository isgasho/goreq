package log

import (
	"fmt"

	"github.com/aiscrm/goreq"
)

func Dump() goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			next(req, resp)
			fmt.Println(resp.Dump())
			return nil
		}
	}
}
