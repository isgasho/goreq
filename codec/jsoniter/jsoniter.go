package jsoniter

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/aiscrm/goreq/codec"
)

func NewCodec(opts ...codec.Option) codec.Codec {
	options := codec.Options{
		EscapeHTML: true,
	}
	c := &jsonCodec{
		opts: options,
	}
	c.Init(opts...)
	return c
}

type jsonCodec struct {
	opts codec.Options
	json jsoniter.API
}

func (j *jsonCodec) Init(options ...codec.Option) {
	for _, o := range options {
		o(&j.opts)
	}
	j.json = jsoniter.Config{
		EscapeHTML:             j.opts.EscapeHTML,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
}

func (j *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}

func (j *jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return j.json.Unmarshal(data, v)
}
