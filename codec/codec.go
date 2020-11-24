package codec

const (
	JSONCodec = "json"
	XMLCodec  = "xml"
)

type Codec interface {
	Init(...Option)
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type Marshal func(interface{}) ([]byte, error)
type Unmarshal func([]byte, interface{}) error

type Options struct {
	EscapeHTML bool
	Indent     string
	Prefix     string
}

type Option func(options *Options)

func DisableEscapeHTML() Option {
	return func(options *Options) {
		options.EscapeHTML = false
	}
}

func EnableEscapeHTML() Option {
	return func(options *Options) {
		options.EscapeHTML = true
	}
}

func SetIndent(indent string) Option {
	return func(options *Options) {
		options.Indent = indent
	}
}

func SetPrefix(prefix string) Option {
	return func(options *Options) {
		options.Prefix = prefix
	}
}
