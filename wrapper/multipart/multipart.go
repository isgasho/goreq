package multipart

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

// File upload file with custom field name and file name
func File(fieldName, fileName string, file io.ReadCloser) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			util.AddMultipartFile(request, fieldName, fileName, file)
			return next(response, request)
		}
	}
}

// FileBytes upload file with custom field name and file name
func FileBytes(fieldName, fileName string, data []byte) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			if len(data) == 0 {
				return next(response, request)
			}
			util.AddMultipartFile(request, fieldName, fileName, ioutil.NopCloser(bytes.NewReader(data)))
			return next(response, request)
		}
	}
}
