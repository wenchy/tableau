package tableau

type Format int

// file format
const (
	JSON      Format = 0
	Protowire        = 1
	Prototext        = 2
)

// Options is the wrapper of tableau params.
type Options struct {
	ProtoPackage              string // protobuf package name.
	InputDir                  string // input dir of workbooks.
	OutputDir                 string // output dir of generated files.
	OutputFilenameAsSnakeCase bool   // output filename as snake case, default is camel case same as the protobuf message name.
	OutputFormat              Format // output format: json, protowire, or prototext, and default is json.
	OutputPretty              bool   // output pretty format, with mulitline and indent.
	LocationName              string // Location represents the collection of time offsets in use in a geographical area. Default is "Asia/Shanghai".
	EmitUnpopulated           bool   // EmitUnpopulated specifies whether to emit unpopulated fields.
	LogLevel                  string // Log level: debug, info, warn, error
}

func (opts *Options) init() {
	if opts.LocationName == "" {
		opts.LocationName = "Asia/Shanghai"
	}
	if opts.LogLevel == "" {
		opts.LogLevel = "info"
	}
}
