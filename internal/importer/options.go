package importer

import (
	"github.com/Wenchy/tableau/proto/tableaupb"
	"google.golang.org/protobuf/proto"
)

type format int

// file format
const (
	Excel format = 0
	CSV   format = 1
	XML   format = 2
)

type SheetParser interface {
	Parse(protomsg proto.Message, sheet *Sheet, wsOpts *tableaupb.WorksheetOptions) error
}

type Options struct {
	Format format      // file format: Excel, CSV, XML. Default: Excel.
	Sheets []string    // sheet names to import
	Parser SheetParser // parser to parse the worksheet
}

// Option is the functional option type.
type Option func(*Options)

func Format(fmt format) Option {
	return func(opts *Options) {
		opts.Format = fmt
	}
}

func Sheets(sheets []string) Option {
	return func(opts *Options) {
		opts.Sheets = sheets
	}
}

func Parser(sp SheetParser) Option {
	return func(opts *Options) {
		opts.Parser = sp
	}
}

func newDefaultOptions() *Options {
	return &Options{
		Format: Excel,
	}
}

func parseOptions(setters ...Option) *Options {
	// Default Options
	opts := newDefaultOptions()
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}
