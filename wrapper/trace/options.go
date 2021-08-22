package trace

type Options struct {
	DumpRequest  bool
	DumpResponse bool
}

type Option func(options *Options)

func DumpRequest(dumpRequest bool) Option {
	return func(options *Options) {
		options.DumpRequest = dumpRequest
	}
}

func DumpResponse(dumpResponse bool) Option {
	return func(options *Options) {
		options.DumpResponse = dumpResponse
	}
}
