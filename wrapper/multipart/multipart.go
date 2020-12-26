package multipart

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq"
)

// File upload file with custom field name and file name
func File(fieldName, fileName string, file io.ReadCloser) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			util.AddMultipartFile(req.Request, fieldName, fileName, file)
			return next(req, resp)
		}
	}
}

// FileBytes upload file with custom field name and file name
func FileBytes(fieldName, fileName string, data []byte) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			if len(data) == 0 {
				return next(req, resp)
			}
			util.AddMultipartFile(req.Request, fieldName, fileName, ioutil.NopCloser(bytes.NewReader(data)))
			return next(req, resp)
		}
	}
}
