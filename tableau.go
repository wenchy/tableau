package tableau

import (
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/converter"
	"github.com/Wenchy/tableau/internal/protogen"
	"github.com/Wenchy/tableau/internal/xlsxgen"
)

// Tableaux is an alias type of converter.Tableaux.
type Tableaux = converter.Tableaux

func Convert(protoPackage, indir, outdir string) {
	tableaux := converter.Tableaux{ProtoPackage: protoPackage, InputDir: indir, OutputDir: outdir}
	tableaux.Convert()
}

// NewTableaux creates a new Tableaux with specified options.
func NewTableaux(opts *Options) *Tableaux {
	opts.init()
	tbx := converter.Tableaux{
		ProtoPackage:              opts.ProtoPackage,
		InputDir:                  opts.InputDir,
		OutputDir:                 opts.OutputDir,
		OutputFilenameAsSnakeCase: opts.OutputFilenameAsSnakeCase,
		OutputFormat:              converter.Format(opts.OutputFormat),
		OutputPretty:              opts.OutputPretty,
		LocationName:              opts.LocationName,
		EmitUnpopulated:           opts.EmitUnpopulated,
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
func Xlsx2Protoconf(protoPackage, goPackage, indir, outdir string) {
	g := protogen.Generator{
		ProtoPackage: protoPackage,
		GoPackage:    goPackage,
		InputDir:     indir,
		OutputDir:    outdir,
	}
	g.Generate()
}
