package vo

import (
	"encoding/json"
	"testing"
)

func TestJSONValue(t *testing.T) {
	name := "This is name"
	age := 23
	v := &JSONValue{}
	if err := json.Unmarshal([]byte(`{"name":"This is name","age":23}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.MustString("name") != name {
		t.Errorf("want %s, but %s", name, v.MustString("name"))
	}
	if v.MustInt("age") != age {
		t.Errorf("want %d, but %d", age, v.MustInt("age"))
	}
}
