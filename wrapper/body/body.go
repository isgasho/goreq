package body

import (
	"bytes"
	"encoding"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/codec"

	"github.com/aiscrm/goreq/util"
)

// JSON convert body to json data
func JSON(body interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderContentType, util.HeaderContentTypeJSON)
			data, err := json.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(response, request)
		}
	}
}

// JSONWithCodec convert body to json data with custom codec
func JSONWithCodec(body interface{}, c codec.Codec) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderContentType, util.HeaderContentTypeJSON)
			data, err := c.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(response, request)
		}
	}
}

// XML convert body to xml data
func XML(body interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			request.Header.Set(util.HeaderContentType, util.HeaderContentTypeXML)
			data, err := xml.Marshal(body)
			if err != nil {
				return err
			}
			return Binary(data)(next)(response, request)
		}
	}
}

// Body with body
func Body(body interface{}) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
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
			return Binary(data)(next)(response, request)
		}
	}
}

// Binary with bytes body
func Binary(body []byte) wrapper.CallWrapper {
	return Reader(bytes.NewReader(body))
}

func Reader(body io.Reader) wrapper.CallWrapper {
	return func(next wrapper.CallFunc) wrapper.CallFunc {
		return func(response *http.Response, request *http.Request) error {
			if err := util.SetBinary(request, body); err != nil {
				return err
			}
			return next(response, request)
		}
	}
}
