package vo

import (
	"encoding/json"
	"errors"

	"github.com/buger/jsonparser"
)

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte{}, nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (m RawMessage) GetString(keys ...string) (string, error) {
	return jsonparser.GetString(m, keys...)
}

func (m RawMessage) MustString(keys ...string) string {
	v, err := m.GetString(keys...)
	if err != nil {
		return ""
	}
	return v
}

func (m RawMessage) GetInt(keys ...string) (int, error) {
	v, err := m.GetInt64(keys...)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (m RawMessage) MustInt(keys ...string) int {
	v, err := m.GetInt(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m RawMessage) GetInt64(keys ...string) (int64, error) {
	return jsonparser.GetInt(m, keys...)
}

func (m RawMessage) MustInt64(keys ...string) int64 {
	v, err := m.GetInt64(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m RawMessage) GetFloat(keys ...string) (float64, error) {
	return jsonparser.GetFloat(m, keys...)
}

func (m RawMessage) MustFloat(keys ...string) float64 {
	v, err := m.GetFloat(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m RawMessage) GetBool(keys ...string) (bool, error) {
	return jsonparser.GetBoolean(m, keys...)
}

func (m RawMessage) MustBool(keys ...string) bool {
	v, err := m.GetBool(keys...)
	if err != nil {
		return false
	}
	return v
}

var _ json.Marshaler = (*RawMessage)(nil)
var _ json.Unmarshaler = (*RawMessage)(nil)
