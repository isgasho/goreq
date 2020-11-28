package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/aiscrm/goreq/codec"
)

type Options struct {
	EnableCookie          bool
	Timeout               time.Duration
	DialTimeout           time.Duration
	DialKeepAlive         time.Duration
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	Transport             http.RoundTripper
	TLSClientConfig       *tls.Config
	Proxy                 func(*http.Request) (*url.URL, error)
	Codec                 codec.Codec
	Errors                []error
}

type Option func(options *Options)

func WithTransport(transport http.RoundTripper) Option {
	return func(options *Options) {
		options.Transport = transport
	}
}

// EnableInsecureTLS allows insecure https
func EnableInsecureTLS(enable bool) Option {
	return func(options *Options) {
		if options.TLSClientConfig == nil {
			options.TLSClientConfig = &tls.Config{}
		}
		options.TLSClientConfig.InsecureSkipVerify = enable
	}
}

// EnableCookie enable or disable cookie manager
func EnableCookie(enable bool) Option {
	return func(options *Options) {
		options.EnableCookie = enable
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithTLSCert(certPEMBlock, keyPEMBlock []byte) Option {
	return func(options *Options) {
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			options.Errors = append(options.Errors, err)
			return
		}
		if options.TLSClientConfig == nil {
			options.TLSClientConfig = &tls.Config{}
		}
		options.TLSClientConfig.Certificates = append(options.TLSClientConfig.Certificates, cert)
	}
}

func WithProxy(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(options *Options) {
		options.Proxy = proxy
	}
}

func WithProxyURL(proxyURL string) Option {
	return func(options *Options) {
		u, err := url.Parse(proxyURL)
		if err != nil {
			options.Errors = append(options.Errors, err)
			return
		}
		options.Proxy = http.ProxyURL(u)
	}
}

func WithCodec(codec codec.Codec) Option {
	return func(options *Options) {
		options.Codec = codec
	}
}

// TODO
//func WithPrometheus() Option {
//	return func(options *Options) {
//		//http.Handle("/metrics", promhttp.Handler())
//	}
//}
