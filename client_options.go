package goreq

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type ClientOptions struct {
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
	Errors                []error
}

type ClientOption func(options *ClientOptions)

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(options *ClientOptions) {
		options.Transport = transport
	}
}

// EnableInsecureTLS allows insecure https
func EnableInsecureTLS(enable bool) ClientOption {
	return func(options *ClientOptions) {
		if options.TLSClientConfig == nil {
			options.TLSClientConfig = &tls.Config{}
		}
		options.TLSClientConfig.InsecureSkipVerify = enable
	}
}

// EnableCookie enable or disable cookie manager
func EnableCookie(enable bool) ClientOption {
	return func(options *ClientOptions) {
		options.EnableCookie = enable
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

func WithTLSCert(certPEMBlock, keyPEMBlock []byte) ClientOption {
	return func(options *ClientOptions) {
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

func WithProxy(proxy func(*http.Request) (*url.URL, error)) ClientOption {
	return func(options *ClientOptions) {
		options.Proxy = proxy
	}
}

func WithProxyURL(proxyURL string) ClientOption {
	return func(options *ClientOptions) {
		u, err := url.Parse(proxyURL)
		if err != nil {
			options.Errors = append(options.Errors, err)
			return
		}
		options.Proxy = http.ProxyURL(u)
	}
}
