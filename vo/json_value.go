package vo

import (
	"encoding/json"
	"errors"

	"github.com/buger/jsonparser"
)

// JSONValue is a raw encoded JSON value.
// It implements Marshaller and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type JSONValue []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m JSONValue) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte{}, nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSONValue) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.JSONValue: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (m JSONValue) GetString(keys ...string) (string, error) {
	return jsonparser.GetString(m, keys...)
}

func (m JSONValue) MustString(keys ...string) string {
	v, err := m.GetString(keys...)
	if err != nil {
		return ""
	}
	return v
}

func (m JSONValue) GetInt(keys ...string) (int, error) {
	v, err := m.GetInt64(keys...)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (m JSONValue) MustInt(keys ...string) int {
	v, err := m.GetInt(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m JSONValue) GetInt64(keys ...string) (int64, error) {
	return jsonparser.GetInt(m, keys...)
}

func (m JSONValue) MustInt64(keys ...string) int64 {
	v, err := m.GetInt64(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m JSONValue) GetFloat(keys ...string) (float64, error) {
	return jsonparser.GetFloat(m, keys...)
}

func (m JSONValue) MustFloat(keys ...string) float64 {
	v, err := m.GetFloat(keys...)
	if err != nil {
		return 0
	}
	return v
}

func (m JSONValue) GetBool(keys ...string) (bool, error) {
	return jsonparser.GetBoolean(m, keys...)
}

func (m JSONValue) MustBool(keys ...string) bool {
	v, err := m.GetBool(keys...)
	if err != nil {
		return false
	}
	return v
}

var _ json.Marshaler = (*JSONValue)(nil)
var _ json.Unmarshaler = (*JSONValue)(nil)
