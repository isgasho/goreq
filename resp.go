package goreq

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq/codec"
)

// Resp represents a http response
type Resp struct {
	Request  *http.Request
	Response *http.Response
	body     []byte
	Error    error
	Cost     time.Duration
	Timeout  bool
	codec    codec.Codec
}

// StatusCode returns status code
func (r *Resp) StatusCode() int {
	return r.Response.StatusCode
}

// ContentLength returns content length
func (r *Resp) ContentLength() int64 {
	return r.Response.ContentLength
}

// ContentType returns content type
func (r *Resp) ContentType() string {
	return r.Response.Header.Get(util.HeaderContentType)
}

// Consume close response body
func (r *Resp) Consume(read bool) {
	if read {
		r.AsBytes()
	} else if r.body == nil {
		r.Response.Body.Close()
		r.body = []byte{}
	}
}

// Bytes returns response body as []byte
func (r *Resp) Bytes() []byte {
	data, _ := r.AsBytes()
	return data
}

// AsBytes returns response body as []byte,
// return error if error happend when reading
// the response body
func (r *Resp) AsBytes() ([]byte, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	if r.body != nil {
		return r.body, nil
	}
	if r.Response.Body == nil {
		return []byte{}, nil
	}
	defer r.Response.Body.Close()
	r.body, r.Error = ioutil.ReadAll(r.Response.Body)
	return r.body, r.Error
}

// AsReader returns response body as reader
func (r *Resp) AsReader() (io.Reader, error) {
	data, err := r.AsBytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

// String returns response body as string
func (r *Resp) String() string {
	data, _ := r.AsBytes()
	return string(data)
}

// AsString returns response body as string,
// return error if error happend when reading
// the response body
func (r *Resp) AsString() (string, error) {
	data, err := r.AsBytes()
	return string(data), err
}

// AsStruct convert to struct. default to use json format
func (r *Resp) AsStruct(v interface{}, unmarshal codec.Unmarshal) error {
	data, err := r.AsBytes()
	if err != nil {
		return err
	}
	return unmarshal(data, v)
}

// AsJSONStruct convert json response body to struct or map
func (r *Resp) AsJSONStruct(v interface{}) error {
	return r.AsStruct(v, r.codec.Unmarshal)
}

// AsXMLStruct convert xml response body to struct or map
func (r *Resp) AsXMLStruct(v interface{}) error {
	return r.AsStruct(v, xml.Unmarshal)
}

func (r *Resp) AsJSONMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := r.AsJSONStruct(&m)
	return m, err
}

func (r *Resp) AsXMLMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := r.AsXMLStruct(&m)
	return m, err
}

// AsFile save to file
func (r *Resp) AsFile(dest string) error {
	data, err := r.AsBytes()
	if err != nil {
		return err
	}
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
	// if r.body != nil {
	// 	_, Error = file.Write(r.body)
	// 	return Error
	// }

	// defer r.response.Body.Close()
	// _, Error = io.Copy(file, r.response.Body)
	// return Error
}

func (r *Resp) Dump() string {
	buf := &bytes.Buffer{}
	if r.Request != nil {
		reqData, _ := httputil.DumpRequest(r.Request, true)
		buf.Write(reqData)
	}
	buf.WriteString("\n\n========================================\n")
	if r.Response != nil {
		respData, _ := httputil.DumpResponse(r.Response, true)
		buf.Write(respData)
	}
	return string(buf.Bytes())
}
