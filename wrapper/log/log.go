package log

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/aiscrm/goreq/wrapper"
)

func Dump() wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			next(response, request)

			buf := &bytes.Buffer{}
			if request != nil {
				reqData, _ := httputil.DumpRequest(request, true)
				buf.Write(reqData)
			}
			buf.WriteString("\n========================================\n")
			if response != nil {
				respData, _ := httputil.DumpResponse(response, true)
				buf.Write(respData)
			}
			fmt.Println(string(buf.Bytes()))
			return nil
		}
	}
}
