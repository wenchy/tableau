package tableau

import (
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen"
	"github.com/Wenchy/tableau/internal/protogen"
	"github.com/Wenchy/tableau/internal/xlsxgen"
	"github.com/Wenchy/tableau/options"
)

// NewTableaux creates a new Tableaux with specified options.
func NewTableaux(protoPackage, indir, outdir string, setters ...options.Option) *confgen.Tableaux {
	opts := options.ParseOptions(setters...)
	tbx := confgen.Tableaux{
		ProtoPackage: protoPackage,
		InputDir:     indir,
		OutputDir:    outdir,

		LocationName: opts.LocationName,

		Output: opts.Output,
	}
	atom.InitZap(opts.LogLevel)
	atom.Log.Infof("options inited: %+v", opts)
	return &tbx
}

// Protoconf2Xlsx converts protoconf files to xlsx files (with meta header).
func Protoconf2Xlsx(protoPackage, indir, outdir string) {
	g := xlsxgen.Generator{
		ProtoPackage: protoPackage,
		InputDir:     indir,
		OutputDir:    outdir,
	}
	g.Generate()
}

// Xlsx2Protoconf converts xlsx files (with meta header) to protoconf files.
func Xlsx2Protoconf(protoPackage, goPackage, indir, outdir string, setters ...options.Option) {
	opts := options.ParseOptions(setters...)
	g := protogen.Generator{
		ProtoPackage: protoPackage,
		GoPackage:    goPackage,
		InputDir:     indir,
		OutputDir:    outdir,
		Header:       opts.Header,
		Imports:      opts.Imports,
	}

	if err := g.Generate(); err != nil {
		atom.Log.Panic(err)
	}
}
