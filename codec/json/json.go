package json

import (
	"bytes"
	"encoding/json"

	"github.com/aiscrm/goreq/codec"
)

func NewCodec(opts ...codec.Option) codec.Codec {
	options := codec.Options{
		EscapeHTML: true,
	}
	c := &stdCodec{
		opts: options,
	}
	c.Init(opts...)
	return c
}

type stdCodec struct {
	opts codec.Options
}

func (j *stdCodec) Init(options ...codec.Option) {
	for _, o := range options {
		o(&j.opts)
	}
}

func (j *stdCodec) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(j.opts.EscapeHTML)
	enc.SetIndent(j.opts.Prefix, j.opts.Indent)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func (j *stdCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
