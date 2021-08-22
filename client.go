package goreq

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"

	"github.com/aiscrm/goreq/wrapper"

	"github.com/aiscrm/goreq/util"
)

var (
	DefaultClient = NewClient()
)

//type HandlerFunc func(ctx *Context)
//type HandlerChain []HandlerFunc

type Client interface {
	Init(...ClientOption) error
	Options() ClientOptions
	Use(...wrapper.CallWrapper) Client
	Do(*Req, ...ClientOption) *Resp
	New() *Req
	Get(rawURL string) *Req
	Post(rawURL string) *Req
}

func NewClient(opts ...ClientOption) Client {
	// default options
	options := ClientOptions{
		EnableCookie:          true,
		Timeout:               0,
		DialTimeout:           30 * time.Second,
		DialKeepAlive:         30 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Transport:             nil,
		TLSClientConfig:       nil,
		Proxy:                 nil,
		Errors:                []error{},
	}
	c := &client{
		opts: options,
	}
	_ = c.Init(opts...)

	return c
}

type client struct {
	opts       ClientOptions
	httpClient *http.Client
	wrappers   []wrapper.CallWrapper
	//handler    CallFunc
	pool sync.Pool
}

func (c *client) Init(opts ...ClientOption) error {
	for _, o := range opts {
		o(&c.opts)
	}
	// init http client
	c.httpClient = newHttpClient(c.opts)
	return nil
}

func (c *client) Options() ClientOptions {
	return c.opts
}

func (c *client) Use(wrappers ...wrapper.CallWrapper) Client {
	c.wrappers = append(c.wrappers, wrappers...)
	return c
	//nc := &client{
	//	opts: c.opts,
	//}
	//nc.httpClient = newHttpClient(nc.opts)
	//nc.wrappers = make([]CallWrapper, 0, len(c.wrappers)+len(wrappers))
	//nc.wrappers = append(nc.wrappers, c.wrappers...)
	//nc.wrappers = append(nc.wrappers, wrappers...)
	//return nc
}

func (c *client) Do(req *Req, opts ...ClientOption) *Resp {
	for _, o := range opts {
		o(&c.opts)
	}
	resp := new(Resp)
	resp.Response = new(http.Response)
	chain := wrapper.New(req.wrappers...)
	if len(c.wrappers) > 0 {
		chain = chain.Append(c.wrappers...)
	}
	before := time.Now()
	err := chain.Then(c.do)(resp.Response, req.Request)
	resp.Cost = time.Now().Sub(before)
	resp.Request = req.Request
	if err != nil {
		resp.Error = err
		if strings.Contains(resp.Error.Error(), "Client.Timeout exceeded") { // 超时的判断
			resp.Timeout = true
		}
		return resp
	}

	if resp.Response.Header.Get(util.HeaderContentEncoding) == util.HeaderContentEncodingGzip {
		body, err := gzip.NewReader(resp.Response.Body)
		if err == nil {
			resp.Response.Body = body
		}
	}
	if resp.Response.Header.Get(util.HeaderContentEncoding) == util.HeaderContentEncodingDeflate {
		body, err := zlib.NewReader(resp.Response.Body)
		if err == nil {
			resp.Response.Body = body
		}
	}
	return resp
}

func (c *client) New() *Req {
	return New().WithClient(c)
}

func (c *client) Get(rawURL string) *Req {
	return Get(rawURL).WithClient(c)
}

func (c *client) Post(rawURL string) *Req {
	return Post(rawURL).WithClient(c)
}

func (c *client) do(response *http.Response, request *http.Request) error {
	var err error
	var reqBody []byte
	if request.Body != nil {
		reqBody, err = ioutil.ReadAll(request.Body)
		if err == nil {
			request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		}
	}
	res, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	if len(reqBody) > 0 {
		request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	}
	*response = *res
	return nil
}

func newHttpClient(options ClientOptions) *http.Client {
	jar, _ := cookiejar.New(nil)
	if !options.EnableCookie {
		jar = nil
	}
	transport := options.Transport
	if transport == nil {
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   options.DialTimeout,
				KeepAlive: options.DialKeepAlive,
				//DualStack: true,
			}).DialContext,
			MaxIdleConns:          options.MaxIdleConns,
			IdleConnTimeout:       options.IdleConnTimeout,
			TLSHandshakeTimeout:   options.TLSHandshakeTimeout,
			TLSClientConfig:       options.TLSClientConfig,
			ExpectContinueTimeout: options.ExpectContinueTimeout,
		}
	}
	return &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   options.Timeout,
	}
}
