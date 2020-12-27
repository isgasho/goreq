package test

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/aiscrm/goreq/wrapper/breaker"

	"github.com/aiscrm/goreq/wrapper/breaker/hystrix"

	"github.com/aiscrm/goreq/codec"
	"github.com/aiscrm/goreq/codec/json"
	"github.com/aiscrm/goreq/vo"

	"github.com/aiscrm/goreq/wrapper/body"

	"github.com/aiscrm/goreq/wrapper/form"
	"github.com/aiscrm/goreq/wrapper/multipart"

	"github.com/aiscrm/goreq/wrapper/log"

	"github.com/aiscrm/goreq"

	"github.com/aiscrm/goreq/wrapper/query"

	"github.com/aiscrm/goreq/wrapper/url"
)

func TestReq_AddQueryParams(t *testing.T) {
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		queries := r.URL.Query()
		for k, v := range queries {
			fmt.Println(k, v)
		}
		fmt.Println(r.Form)
		fmt.Println(r.PostForm)
		title := r.FormValue("title")
		//f, fh, _ := r.FormFile("media")
		//fmt.Println("filename=", fh.Filename, ", size=", fh.Size)
		//d, _ := ioutil.ReadAll(f)
		//fmt.Println("文件内容:", string(d))
		//name := r.URL.Query().Get("name")
		w.Write([]byte(title))
	}
	ts := httptest.NewServer(http.HandlerFunc(getHandler))

	data, err := goreq.NewClient().Use(log.Dump()).Post(ts.URL).
		Use(
			url.Path("/hello"),
			query.SetMap(map[string]interface{}{
				"name": "corel",
				"age":  18,
			}),
			form.AddMap(map[string]interface{}{
				"book": "This is a book,哈哈",
			}),
			multipart.FileBytes("media", "1", []byte(`你好，world`)),
			form.Set("title", "ThisIsTitle"),
			multipart.FileBytes("image", "2", []byte(`你好，图片`)),
			//log.Dump(),
		).Do().AsString()
	if err != nil {
		t.Error(err)
	}
	if data != "ThisIsTitle" {
		t.Errorf("want %s, but %s", "ThisIsTitle", data)
	}
}

func TestBody(t *testing.T) {
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		w.Write(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(getHandler))

	type Scene struct {
		SceneID string `json:"scene_id"`
	}

	type ActionInfo struct {
		Scene Scene `json:"scene"`
	}

	type QrCodeRequest struct {
		ActionName string     `json:"action_name"`
		ActionInfo ActionInfo `json:"action_info"`
	}
	qrCodeRequest := QrCodeRequest{
		ActionName: "QR_LIMIT_SCENE",
		ActionInfo: ActionInfo{
			Scene: Scene{SceneID: "id=1&name=corel"},
		},
	}
	qrCodeResponse := QrCodeRequest{}
	err := goreq.NewClient().Use(log.Dump()).Post(ts.URL).
		Use(
			body.JSONWithCodec(qrCodeRequest, json.NewCodec(codec.DisableEscapeHTML())),
		).Do().AsJSONStruct(&qrCodeResponse)
	if err != nil {
		t.Error(err)
	}
	if qrCodeResponse.ActionInfo.Scene.SceneID != qrCodeRequest.ActionInfo.Scene.SceneID {
		t.Errorf("want %s, but %s", qrCodeRequest.ActionInfo.Scene.SceneID, qrCodeResponse.ActionInfo.Scene.SceneID)
	}
}

func TestJSONValue(t *testing.T) {
	name := "你好"
	age := 23
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		rname := r.Form.Get("name")
		rage, _ := strconv.Atoi(r.Form.Get("age"))
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		p := Person{
			Name: rname,
			Age:  rage,
		}
		data, _ := stdjson.Marshal(p)
		w.Write(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(getHandler))

	v := vo.JSONValue{}
	err := goreq.NewClient().Use(log.Dump()).Post(ts.URL).
		Use(
			query.Set("name", name),
			query.Set("age", age),
		).
		Do().AsJSONStruct(&v)
	if err != nil {
		t.Error(err)
	}
	if v.MustString("name") != name {
		t.Errorf("want %s, but %s", name, v.MustString("name"))
	}
	if v.MustInt("age") != age {
		t.Errorf("want %d, but %d", age, v.MustInt("age"))
	}
}

func TestBreaker(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		sleep := r.URL.Query().Get("sleep")
		sleepMs, _ := strconv.Atoi(sleep)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		w.Write([]byte(sleep))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))

	sleep := 2000
	client := goreq.NewClient().Use(hystrix.Breaker(hystrix.DefaultTimeout(3000)))
	resp := client.New().WithMethod(http.MethodPost).
		Use(
			url.URL(ts.URL),
			query.AddMap(map[string]interface{}{
				"sleep": sleep,
			}),
		).Do()
	if resp.Error != nil {
		t.Errorf("breaker want to be closed, but is open: %v", resp.Error)
	}

	sleep = 8000
	resp2 := client.New().WithMethod(http.MethodPost).
		Use(
			url.URL(ts.URL),
			query.AddMap(map[string]interface{}{
				"sleep": sleep,
			}),
		).Do()
	if resp2.Error == nil {
		t.Errorf("breaker want to be open, but is closed")
	}
	if errors.Is(resp2.Error, breaker.CircuitError{}) {
		t.Errorf("is not CircuitError: %v", resp2.Error)
	}
}

func BenchmarkGet(b *testing.B) {
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}
	ts := httptest.NewServer(http.HandlerFunc(getHandler))

	for i := 0; i < b.N; i++ {
		_, _ = goreq.New().WithMethod(http.MethodPost).
			Use(
				url.URL(ts.URL),
				query.AddMap(map[string]interface{}{
					"name": "corel",
					"age":  18,
				}),
			).Do().AsString()
	}
}
