package body

import (
	"bytes"
	"encoding"
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/aiscrm/goreq/codec"

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq"
)

// JSON convert body to json data
func JSON(body interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderContentType, util.HeaderContentTypeJSON)
			data, err := json.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(req, resp)
		}
	}
}

// JSONWithCodec convert body to json data with custom codec
func JSONWithCodec(body interface{}, c codec.Codec) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderContentType, util.HeaderContentTypeJSON)
			data, err := c.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(req, resp)
		}
	}
}

// XML convert body to xml data
func XML(body interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			req.Request.Header.Set(util.HeaderContentType, util.HeaderContentTypeXML)
			data, err := xml.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(req, resp)
		}
	}
}

// Body with body
func Body(body interface{}) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			var data []byte
			var err error
			switch b := body.(type) {
			case json.Marshaler:
				data, err = b.MarshalJSON()
			case encoding.BinaryMarshaler:
				data, err = b.MarshalBinary()
			case io.ReadCloser:
				var buf bytes.Buffer
				if _, err := buf.ReadFrom(b); err != nil {
					return err
				}
				if err := b.Close(); err != nil {
					return err
				}
				data = buf.Bytes()
			case io.Reader:
				var buf bytes.Buffer
				if _, err := buf.ReadFrom(b); err != nil {
					return err
				}
				data = buf.Bytes()
			case bytes.Buffer:
				data = b.Bytes()
			case string:
				data = []byte(b)
			case []byte:
				data = b
			case func() ([]byte, error):
				data, err = b()
			default:
				// or return error
				data, err = json.Marshal(body)
			}
			if err != nil {
				return err
			}
			return Binary(data)(next)(req, resp)
		}
	}
}

// Body with body
func Binary(body []byte) goreq.CallWrapper {
	return Reader(bytes.NewReader(body))
}

func Reader(body io.Reader) goreq.CallWrapper {
	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			if err := util.SetBinary(req.Request, body); err != nil {
				return err
			}
			return next(req, resp)
		}
	}
}
