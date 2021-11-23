package options

type Format int

// file format
const (
	JSON      Format = 0
	Protowire        = 1
	Prototext        = 2
)

// Options is the wrapper of tableau params.
// Options follow the design of Functional Options (https://github.com/tmrts/go-patterns/blob/master/idiom/functional-options.md).
type Options struct {
	LocationName string        // Location represents the collection of time offsets in use in a geographical area. Default is "Asia/Shanghai".
	LogLevel     string        // Log level: debug, info, warn, error
	Header       *HeaderOption // header rows of excel file.
	Output       *OutputOption // output settings.
	Imports      []string      // imported common proto file paths
}

type HeaderOption struct {
	Namerow int32
	Typerow int32
	Noterow int32
	Datarow int32
}

type OutputOption struct {
	FilenameAsSnakeCase bool   // output filename as snake case, default is camel case same as the protobuf message name.
	Format              Format // output pretty format, with mulitline and indent.
	Pretty              bool   // output format: json, protowire, or prototext, and default is json.
	// Output.EmitUnpopulated specifies whether to emit unpopulated fields. It does not
	// emit unpopulated oneof fields or unpopulated extension fields.
	// The JSON value emitted for unpopulated fields are as follows:
	//  ╔═══════╤════════════════════════════╗
	//  ║ JSON  │ Protobuf field             ║
	//  ╠═══════╪════════════════════════════╣
	//  ║ false │ proto3 boolean fields      ║
	//  ║ 0     │ proto3 numeric fields      ║
	//  ║ ""    │ proto3 string/bytes fields ║
	//  ║ null  │ proto2 scalar fields       ║
	//  ║ null  │ message fields             ║
	//  ║ []    │ list fields                ║
	//  ║ {}    │ map fields                 ║
	//  ╚═══════╧════════════════════════════╝
	EmitUnpopulated bool
}

// Option is the functional option type.
type Option func(*Options)

func Header(o *HeaderOption) Option {
	return func(opts *Options) {
		opts.Header = o
	}
}

func Output(o *OutputOption) Option {
	return func(opts *Options) {
		opts.Output = o
	}
}

func LocationName(o string) Option {
	return func(opts *Options) {
		opts.LocationName = o
	}
}

func LogLevel(level string) Option {
	return func(opts *Options) {
		opts.LogLevel = level
	}
}

func Imports(imports []string) Option {
	return func(opts *Options) {
		opts.Imports = imports
	}
}

func newDefaultOptions() *Options {
	return &Options{
		LocationName: "Asia/Shanghai",
		LogLevel:     "info",

		Header: &HeaderOption{
			Namerow: 1,
			Typerow: 2,
			Noterow: 3,
			Datarow: 4,
		},
		Output: &OutputOption{
			FilenameAsSnakeCase: false,
			Format:              JSON,
			Pretty:              true,
			EmitUnpopulated:     true,
		},
	}
}

func ParseOptions(setters ...Option) *Options {
	// Default Options
	opts := newDefaultOptions()
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}
