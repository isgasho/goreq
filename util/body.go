package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func SetBinary(req *http.Request, body io.Reader) error {
	if body == nil {
		return nil
	}
	// see http.NewRequestWithContext(ctx context.Context, method, url string, body io.Reader)
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
		}
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
	return nil
}

func ToMultipart(req *http.Request, bodyWriter *multipart.Writer) error {
	postForm := req.PostForm
	if postForm != nil {
		for key, value := range postForm {
			for _, val := range value {
				_ = bodyWriter.WriteField(key, ToString(val))
			}
		}
	}
	req.PostForm = make(url.Values)
	req.ParseMultipartForm(defaultMaxMemory)
	if req.MultipartForm != nil {
		//for key, value := range request.MultipartForm.Value {
		//	for _, val := range value {
		//		_ = bodyWriter.WriteField(key, ToString(val))
		//	}
		//}
		for key, value := range req.MultipartForm.File {
			for _, val := range value {
				fileWriter, err := bodyWriter.CreateFormFile(key, val.Filename)
				if err != nil {
					return err
				}
				f, err := val.Open()
				if err != nil {
					return err
				}
				_, err = io.Copy(fileWriter, f)
				if err != nil {
					return err
				}
			}
		}
	}

	req.PostForm = postForm
	req.Header.Set(HeaderContentType, bodyWriter.FormDataContentType())
	return nil
}

func AddMultipartFile(req *http.Request, fieldName, fileName string, file io.ReadCloser) error {
	data := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(data)

	if err := ToMultipart(req, bodyWriter); err != nil {
		return err
	}

	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	_ = bodyWriter.Close()
	return SetBinary(req, bytes.NewReader(data.Bytes()))
}
