package tableau

import (
	"fmt"
)

type Format int

// file format
const (
	JSON      Format = 0
	Protobin         = 1
	Prototext        = 2
)

// Options is the wrapper of tableau params.
type Options struct {
	ProtoPackageName          string // protobuf package name.
	InputPath                 string // root dir of workbooks.
	OutputPath                string // output path of generated files.
	OutputFilenameAsSnakeCase bool   // output filename as snake case, default is camel case same as the protobuf message name.
	OutputFormat              Format // output format: json, protobin, or prototext, and default is json.
	OutputPretty              bool   // output pretty format, with mulitline and indent.
}

func (opts *Options) init() {
	fmt.Println("options inited")
}
