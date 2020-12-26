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

	"github.com/aiscrm/goreq/util"

	"github.com/aiscrm/goreq/client"
)

var (
	DefaultClient = NewClient()
)

//type HandlerFunc func(ctx *Context)
//type HandlerChain []HandlerFunc

type Client interface {
	Init(...client.Option) error
	Options() client.Options
	Use(...CallWrapper) Client
	Do(*Req, ...client.Option) *Resp
	New() *Req
	Get(rawURL string) *Req
	Post(rawURL string) *Req
}

func NewClient(opts ...client.Option) Client {
	// default options
	options := client.Options{
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
	c := &cli{
		opts: options,
	}
	c.Init(opts...)

	return c
}

type cli struct {
	opts       client.Options
	httpClient *http.Client
	wrappers   []CallWrapper
	//handler    CallFunc
	pool sync.Pool
}

func (c *cli) Init(opts ...client.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	// init http client
	c.httpClient = newHttpClient(c.opts)
	return nil
}

func (c *cli) Options() client.Options {
	return c.opts
}

func (c *cli) Use(wrappers ...CallWrapper) Client {
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

func (c *cli) Do(r *Req, opts ...client.Option) *Resp {
	for _, o := range opts {
		o(&c.opts)
	}
	resp := new(Resp)
	chain := newChain(r.wrappers...)
	if len(c.wrappers) > 0 {
		chain = chain.Append(c.wrappers...)
	}
	err := chain.Then(c.do)(r, resp)
	if err != nil {
		resp.Error = err
	}
	return resp
}

func (c *cli) New() *Req {
	return New().WithClient(c)
}

func (c *cli) Get(rawURL string) *Req {
	return Get(rawURL).WithClient(c)
}

func (c *cli) Post(rawURL string) *Req {
	return Post(rawURL).WithClient(c)
}

func (c *cli) do(req *Req, resp *Resp) error {
	if req.Error != nil {
		return req.Error
	}

	before := time.Now()
	reqBody := req.GetBody()
	resp.Response, resp.Error = c.httpClient.Do(req.Request)
	resp.Cost = time.Now().Sub(before)

	if resp.Error != nil && strings.Contains(resp.Error.Error(), "Client.Timeout exceeded") { // 超时的判断
		resp.Timeout = true
	}
	if resp.Error == nil {
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
	}

	// for dump
	req.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	resp.Request = req.Request

	return nil
}

func newHttpClient(options client.Options) *http.Client {
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
