package hystrix

import (
	"net/http"

	"github.com/aiscrm/goreq"

	"github.com/afex/hystrix-go/hystrix"
)

type Options struct {
	KeyFunc            func(request *http.Request) string
	Commands           map[string]hystrix.CommandConfig
	DefaultCommandName string
}
type Option func(Options)

func KeyFunc(keyFunc func(request *http.Request) string) Option {
	return func(options Options) {
		options.KeyFunc = keyFunc
	}
}

// 如果不想使用默认熔断器，可以设置name为空字符串
// 如果想改变默认熔断器规则，可以使用ConfigureCommand来配置name=ConfigureCommand的熔断器规则
func DefaultCommandName(name string) Option {
	return func(options Options) {
		options.DefaultCommandName = name
	}
}

func ConfigureCommand(name string, config hystrix.CommandConfig) Option {
	return func(options Options) {
		options.Commands[name] = config
		hystrix.ConfigureCommand(name, config)
	}
}

var defaultOptions = Options{
	KeyFunc: func(request *http.Request) string {
		return request.RequestURI
	},
	DefaultCommandName: "default",
}

func Breaker(opts ...Option) goreq.CallWrapper {
	options := defaultOptions
	for _, opt := range opts {
		opt(options)
	}
	if options.DefaultCommandName != "" { // 配置默认的熔断器规则
		if _, ok := options.Commands[options.DefaultCommandName]; !ok {
			defaultConfig := hystrix.CommandConfig{
				Timeout:               1000,
				MaxConcurrentRequests: 100,
				ErrorPercentThreshold: 25,
			}
			options.Commands[options.DefaultCommandName] = defaultConfig
			hystrix.ConfigureCommand(options.DefaultCommandName, defaultConfig)
		}
	}

	return func(next goreq.CallFunc) goreq.CallFunc {
		return func(req *goreq.Req, resp *goreq.Resp) error {
			commandName := options.KeyFunc(req.Request)
			if _, ok := options.Commands[commandName]; !ok && options.DefaultCommandName != "" {
				commandName = options.DefaultCommandName
			}
			hystrix.Do(commandName, func() error {
				return next(req, resp)
			}, func(err error) error {
				return err
			})
			return nil
		}
	}
}
