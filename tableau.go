package tableau

import (
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/converter"
	"github.com/Wenchy/tableau/internal/xlsxgen"
)

// Tableaux is an alias type of converter.Tableaux.
type Tableaux = converter.Tableaux

func Convert(protoPackageName, indir, outdir string) {
	tableaux := converter.Tableaux{ProtoPackageName: protoPackageName, InputDir: indir, OutputDir: outdir}
	tableaux.Convert()
}

func NewTableaux(opts *Options) *Tableaux {
	opts.init()
	tbx := converter.Tableaux{
		ProtoPackageName:          opts.ProtoPackageName,
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

// Generator is an alias type of generator.Generator.
type Generator = xlsxgen.Generator

func Generate(protoPackageName, indir, outdir string) {
	generator := xlsxgen.Generator{ProtoPackageName: protoPackageName, InputDir: indir, OutputDir: outdir}
	generator.Generate()
}
