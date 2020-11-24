package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
