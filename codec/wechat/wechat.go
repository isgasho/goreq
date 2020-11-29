package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aiscrm/goreq/codec"

	"github.com/buger/jsonparser"
)

func NewCodec(opts ...codec.Option) codec.Codec {
	options := codec.Options{
		EscapeHTML: true,
	}
	c := &wechatCodec{
		opts: options,
	}
	c.Init(opts...)
	return c
}

type wechatCodec struct {
	opts codec.Options
}

func (j *wechatCodec) Init(options ...codec.Option) {
	for _, o := range options {
		o(&j.opts)
	}
}

func (j *wechatCodec) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(j.opts.EscapeHTML)
	enc.SetIndent(j.opts.Prefix, j.opts.Indent)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func (j *wechatCodec) Unmarshal(data []byte, v interface{}) error {
	code, _ := jsonparser.GetInt(data, "errcode")
	if code != 0 {
		message, _ := jsonparser.GetString(data, "errmsg")
		return fmt.Errorf(fmt.Sprintf("errcode=%d, errmsg=%s", code, message))
	}
	return json.Unmarshal(data, v)
}
