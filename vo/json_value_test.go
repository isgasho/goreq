package vo

import (
	"encoding/json"
	"testing"
)

func TestJSONValue(t *testing.T) {
	person := struct {
		Name   string  `json:"name"`
		Age    int     `json:"age"`
		Score  float64 `json:"score"`
		Active bool    `json:"active"`
	}{
		Name:   "This is name",
		Age:    23,
		Score:  97.5,
		Active: true,
	}
	data, _ := json.Marshal(person)
	v := &JSONValue{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	if v.MustString("name") != person.Name {
		t.Errorf("want %s, but %s", person.Name, v.MustString("name"))
	}
	if v.MustInt("age") != person.Age {
		t.Errorf("want %d, but %d", person.Age, v.MustInt("age"))
	}
	if v.MustFloat("score") != person.Score {
		t.Errorf("want %v, but %v", person.Score, v.MustFloat("score"))
	}
	if v.MustBool("active") != person.Active {
		t.Errorf("want %v, but %v", person.Active, v.MustBool("active"))
	}
}
